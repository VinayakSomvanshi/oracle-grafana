package plugin

import (
	"encoding/json"
	"testing"
)

func TestParseDatasourceSettings(t *testing.T) {
	t.Run("standard settings with service", func(t *testing.T) {
		raw := json.RawMessage(`{
			"o_hostname": "db.internal",
			"o_port": 1521,
			"o_service": "FREEPDB1",
			"o_user": "grafana_user"
		}`)
		decrypted := map[string]string{
			"o_password": "secret_db_password",
		}

		s := ParseDatasourceSettings(raw, decrypted)
		if s.O_hostname != "db.internal" {
			t.Errorf("expected hostname db.internal, got %s", s.O_hostname)
		}
		if s.O_port != 1521 {
			t.Errorf("expected port 1521, got %d", s.O_port)
		}
		if s.O_service != "FREEPDB1" {
			t.Errorf("expected service FREEPDB1, got %s", s.O_service)
		}
		if s.O_user != "grafana_user" {
			t.Errorf("expected user grafana_user, got %s", s.O_user)
		}
		if s.O_password != "secret_db_password" {
			t.Errorf("expected password secret_db_password, got %s", s.O_password)
		}
		if !s.O_tlsVerify {
			t.Errorf("expected default O_tlsVerify to be true")
		}
	})

	t.Run("cloud TLS and wallet settings", func(t *testing.T) {
		raw := json.RawMessage(`{
			"o_hostname": "adb.us-ashburn-1.oraclecloud.com",
			"o_port": 1522,
			"o_service": "adw_high",
			"o_user": "admin",
			"o_tls": true,
			"o_tlsVerify": false,
			"o_walletPath": "/etc/oracle/wallets/adw"
		}`)
		decrypted := map[string]string{
			"o_password":       "db_pass_123",
			"o_walletPassword": "wallet_secret_456",
		}

		s := ParseDatasourceSettings(raw, decrypted)
		if !s.O_tls {
			t.Errorf("expected O_tls to be true")
		}
		if s.O_tlsVerify {
			t.Errorf("expected O_tlsVerify to be false")
		}
		if s.O_walletPath != "/etc/oracle/wallets/adw" {
			t.Errorf("expected walletPath /etc/oracle/wallets/adw, got %s", s.O_walletPath)
		}
		if s.O_walletPassword != "wallet_secret_456" {
			t.Errorf("expected walletPassword wallet_secret_456, got %s", s.O_walletPassword)
		}
	})

	t.Run("empty and nil parameters", func(t *testing.T) {
		s := ParseDatasourceSettings(nil, nil)
		if !s.O_tlsVerify {
			t.Errorf("expected default O_tlsVerify to be true even with nil inputs")
		}
		if s.O_hostname != "" || s.O_password != "" || s.O_walletPassword != "" {
			t.Errorf("expected empty strings for unset parameters")
		}
	})
}
