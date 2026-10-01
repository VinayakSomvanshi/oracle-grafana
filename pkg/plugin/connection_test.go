package plugin

import (
	"context"
	"strings"
	"testing"

	go_ora "github.com/sijms/go-ora/v2"
)

func TestConnectionLifecycle(t *testing.T) {
	conn := &OracleDatasourceConnection{}

	if conn.IsConnected() {
		t.Fatalf("expected IsConnected to be false initially")
	}

	err := conn.Disconnect()
	if err != nil {
		t.Fatalf("expected disconnect on nil connection to be clean, got: %v", err)
	}

	ctx := context.Background()
	_, err = conn.Conn(ctx)
	if err == nil {
		t.Fatalf("expected error requesting Conn on closed connection")
	}

	err = conn.PingContext(ctx)
	if err == nil {
		t.Fatalf("expected error on PingContext on closed connection")
	}
}

func TestBuildUrlAndCredentialMasking(t *testing.T) {
	settings := &OracleDatasourceSettings{
		O_hostname:       "oracle.internal",
		O_port:           1521,
		O_service:        "FREEPDB1",
		O_user:           "grafana_user",
		O_password:       "P@ssw0rd!SuperSecret",
		O_tls:            true,
		O_tlsVerify:      true,
		O_walletPath:     "/etc/wallets/adw",
		O_walletPassword: "WalletSecretKey99",
	}

	urlOptions := map[string]string{}
	if settings.O_tls {
		urlOptions["SSL"] = "TRUE"
		if settings.O_tlsVerify {
			urlOptions["SSL VERIFY"] = "TRUE"
		} else {
			urlOptions["SSL VERIFY"] = "FALSE"
		}
	}
	if len(settings.O_walletPath) > 0 {
		urlOptions["WALLET"] = settings.O_walletPath
		if len(settings.O_walletPassword) > 0 {
			urlOptions["WALLET PASSWORD"] = settings.O_walletPassword
		}
	}

	connStr := go_ora.BuildUrl(settings.O_hostname, settings.O_port, settings.O_service, settings.O_user, settings.O_password, urlOptions)

	// Verify URL contains expected parameters
	if !strings.Contains(connStr, "SSL=TRUE") {
		t.Errorf("expected connStr to contain SSL=TRUE, got: %s", connStr)
	}
	if !strings.Contains(connStr, "SSL+VERIFY=TRUE") && !strings.Contains(connStr, "SSL VERIFY=TRUE") {
		t.Errorf("expected connStr to contain SSL VERIFY=TRUE, got: %s", connStr)
	}
	if !strings.Contains(connStr, "WALLET=") {
		t.Errorf("expected connStr to contain WALLET option, got: %s", connStr)
	}

	// Test masking logic
	masked := connStr
	if len(settings.O_password) > 0 {
		masked = strings.Replace(masked, settings.O_password, "********", 1)
	}
	if len(settings.O_walletPassword) > 0 {
		masked = strings.Replace(masked, settings.O_walletPassword, "********", 1)
	}

	if strings.Contains(masked, settings.O_password) {
		t.Errorf("CRITICAL SECURITY: Masked connection string leaked database password: %s", masked)
	}
	if strings.Contains(masked, settings.O_walletPassword) {
		t.Errorf("CRITICAL SECURITY: Masked connection string leaked wallet password: %s", masked)
	}
	if !strings.Contains(masked, "********") {
		t.Errorf("Expected masked string to contain asterisks: %s", masked)
	}
}
