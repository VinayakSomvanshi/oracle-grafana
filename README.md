# Hardened Oracle Database DataSource for Grafana

[![CI & Release Pipeline](https://github.com/VinayakSomvanshi/oracle-grafana/actions/workflows/ci.yml/badge.svg)](https://github.com/VinayakSomvanshi/oracle-grafana/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/VinayakSomvanshi/oracle-grafana)](https://goreportcard.com/report/github.com/VinayakSomvanshi/oracle-grafana)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

An enterprise-ready, zero-configuration, strictly read-only **Oracle Database Data Source for Grafana OSS**.

This plugin connects Grafana directly to any Oracle Database (19c, 21c, 23ai, 26ai) with zero database-side configuration.

---

## Zero Database Configuration

You do **not** need DBA access, special privileges, or database-side scripts.

1. Install the plugin.
2. In Grafana, provide your existing Oracle host, port, service name, username, and password.
3. Click **Save & Test**.

The plugin automatically locks every query execution session to read-only mode (`ALTER SESSION SET TRANSACTION READ ONLY`). Even if you connect using a user account that has write or administrative permissions, the database engine itself rejects any write or mutation attempt.

---

## Key Features & Production Hardening

- **Automatic Read-Only Enforcement**: Two-tier protection engine:
  - In-plugin SQL lexer strips comments and enforces that queries start strictly with SELECT or WITH.
  - Automatic session transaction locking ensures the Oracle database kernel blocks any modification attempt (ORA-01456).
- **Pure-Go Driver**: Built on sijms/go-ora/v2 with CGO_ENABLED=0. Requires no Oracle Instant Client RPMs or shared C library dependencies.
- **Native Data Types**: Converts Oracle types (NUMBER, DATE, TIMESTAMP, VARCHAR2) to native Grafana typed DataFrame vectors for instant time-series rendering.
- **Circuit Breaker**: Enforces a 50,000 max row limit with frame warning notices to prevent Out-Of-Memory (OOM) crashes on Grafana.
- **Context Cancellation**: Full propagation of Grafana cancellation contexts to Oracle. When a panel times out or a tab closes, the query terminates in Oracle immediately.
- **Cross-Platform**: Pre-compiled for Linux (amd64, arm64), macOS (Intel, Apple Silicon M1-M4), and Windows.
- **Zero CVEs**: Audited and verified with govulncheck.

---

## Installation & Setup

### 1. Build or Package (or download pre-compiled release)
```bash
make dist
```

### 2. Install on Grafana Host
```bash
sudo ./install_production.sh
```

### 3. Add Data Source in Grafana
1. Navigate to **Connections > Data Sources > Add data source**.
2. Select **Oracle**.
3. Fill in:
   - **Host**: Your Oracle DB IP or hostname (e.g. localhost)
   - **Port**: 1521
   - **Service / SID**: e.g. FREEPDB1 or ORCL
   - **User**: Any standard Oracle user
   - **Password**: User password
4. Click **Save & Test**.

---

## Signing the Plugin (Optional / Production)

To eliminate the need for allow_loading_unsigned_plugins:
1. Create a free organization account on [Grafana.com](https://grafana.com).
2. Generate an Access Policy Token with the plugins:write scope.
3. Sign the plugin:
   ```bash
   export GRAFANA_ACCESS_POLICY_TOKEN="<your-token>"
   export GRAFANA_ROOT_URLS="https://grafana.mycompany.com"
   make sign
   ```

---

## Development & Testing

```bash
# Run unit tests
make test

# Cross-compile all binaries
make build-all-backends

# Clean build artifacts
make clean
```

---

## License

Apache License 2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE) for details.
