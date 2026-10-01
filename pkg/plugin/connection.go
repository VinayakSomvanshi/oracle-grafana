package plugin

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	go_ora "github.com/sijms/go-ora/v2"
)

type OracleDatasourceConnection struct {
	connection *sql.DB
}

func (c *OracleDatasourceConnection) Connect(settings *OracleDatasourceSettings) error {
	var connectionString string
	var err error
	if !c.IsConnected() {
		urlOptions := map[string]string{}
		if len(settings.O_sid) > 0 {
			urlOptions["SID"] = settings.O_sid
		}

		if len(settings.O_connStr) > 0 {
			connectionString = go_ora.BuildJDBC(settings.O_user, settings.O_password, settings.O_connStr, urlOptions)
		} else {
			connectionString = go_ora.BuildUrl(settings.O_hostname, settings.O_port, settings.O_service, settings.O_user, settings.O_password, urlOptions)
		}

		masked := connectionString
		if len(settings.O_password) > 0 {
			masked = strings.Replace(connectionString, settings.O_password, "********", 1)
		}
		log.DefaultLogger.Debug("Connecting to Oracle:", "connStr", masked)

		connection, conErr := sql.Open("oracle", connectionString)
		if conErr != nil {
			log.DefaultLogger.Error("Error opening Oracle connection: ", conErr)
			err = conErr
		} else {
			// Production connection pooling
			connection.SetMaxOpenConns(25)
			connection.SetMaxIdleConns(5)
			connection.SetConnMaxLifetime(30 * time.Minute)
			connection.SetConnMaxIdleTime(5 * time.Minute)

			c.connection = connection
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = c.PingContext(ctx)
		}
	}
	return err
}

func (c *OracleDatasourceConnection) Disconnect() error {
	var err error
	if c.IsConnected() {
		err = c.connection.Close()
		if err != nil {
			log.DefaultLogger.Error("Error closing Oracle connection: ", err)
		}
		c.connection = nil
	}
	return err
}

func (c *OracleDatasourceConnection) IsConnected() bool {
	return c.connection != nil
}

func (c *OracleDatasourceConnection) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.PingContext(ctx)
}

func (c *OracleDatasourceConnection) PingContext(ctx context.Context) error {
	if c.connection != nil {
		return c.connection.PingContext(ctx)
	}
	return fmt.Errorf("oracle connection is closed")
}

func (c *OracleDatasourceConnection) Conn(ctx context.Context) (*sql.Conn, error) {
	if c.connection != nil {
		return c.connection.Conn(ctx)
	}
	return nil, fmt.Errorf("oracle connection is closed")
}

func (c *OracleDatasourceConnection) Reconnect(settings *OracleDatasourceSettings) error {
	if c.IsConnected() {
		_ = c.Disconnect()
	}
	return c.Connect(settings)
}
