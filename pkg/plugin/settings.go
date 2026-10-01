package plugin

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type OracleDatasourceSettings struct {
	O_connStr        string `json:"o_connStr"`
	O_hostname       string `json:"o_hostname"`
	O_password       string `json:"-"`
	O_port           int    `json:"o_port"`
	O_service        string `json:"o_service"`
	O_sid            string `json:"o_sid"`
	O_user           string `json:"o_user"`
	O_tls            bool   `json:"o_tls"`
	O_tlsVerify      bool   `json:"o_tlsVerify"`
	O_walletPath     string `json:"o_walletPath"`
	O_walletPassword string `json:"-"`
}

func ParseDatasourceSettings(rawOptions json.RawMessage, decryptedOptions map[string]string) OracleDatasourceSettings {
	settings := OracleDatasourceSettings{
		O_tlsVerify: true, // Default to true for secure TLS verification
	}
	if decryptedOptions != nil {
		settings.O_password = decryptedOptions["o_password"]
		settings.O_walletPassword = decryptedOptions["o_walletPassword"]
	}
	err := json.Unmarshal(rawOptions, &settings)
	if err != nil {
		log.DefaultLogger.Error("Error parsing Oracle datasource settings: ", err)
	}
	return settings
}
