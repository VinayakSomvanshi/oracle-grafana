# Hardened Oracle Database DataSource for Grafana

[![CI & Release Pipeline](https://github.com/VinayakSomvanshi/oracle-grafana/actions/workflows/ci.yml/badge.svg)](https://github.com/VinayakSomvanshi/oracle-grafana/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/VinayakSomvanshi/oracle-grafana)](https://goreportcard.com/report/github.com/VinayakSomvanshi/oracle-grafana)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

An enterprise-ready, strictly read-only **Oracle Database Data Source for Grafana OSS**.

This project modernizes and hardens the community Oracle Grafana plugin for modern Grafana (v10+ / v11+) and Oracle 19c, 21c, 23ai, and 26ai.

---

## Key Features & Production Hardening

- 🛡️ **Guaranteed Read-Only**: Enforces a 3-tier defense-in-depth model (Lexer sanitization + `ALTER SESSION SET TRANSACTION READ ONLY` + strict DB user permissions). Zero data mutation risk.
- ⚡ **Pure-Go Driver**: Uses `sijms/go-ora/v2` with `CGO_ENABLED=0`. Zero Oracle Instant Client RPMs or shared C library dependencies needed.
- 🏎️ **Native Data Types**: Direct conversion of Oracle types (`NUMBER`, `DATE`, `TIMESTAMP`, `VARCHAR2`) to native Grafana typed DataFrame vectors for instant time-series graph rendering.
- 🛑 **Circuit Breaker**: Enforces a 50,000 max row limit with frame warning notices to prevent Out-Of-Memory (OOM) crashes on Grafana instances.
- ⏱️ **Context Cancellation**: Full propagation of Grafana query cancellation contexts to Oracle. When a panel times out or a tab closes, the query terminates in Oracle immediately.
- 📦 **Cross-Platform**: Pre-compiled for Linux (amd64, arm64), macOS (Intel, Apple Silicon M1-M4), and Windows.
- 🔒 **Zero CVEs**: Audited and verified with `govulncheck`.

---

## Quick Start & Installation

### 1. Configure Oracle Database
Run the hardened SQL script as SYSDBA or SYSTEM on your target PDB (e.g. `FREEPDB1`):
```bash
sqlplus system@localhost:1521/FREEPDB1 @oracle/setup_production.sql
```

### 2. Build or Distribute
```bash
# Build all backend architectures into dist/
make dist
```

### 3. Install on Grafana Host
```bash
sudo ./install_production.sh
```

---

## Signing the Plugin (Optional / Production)

To eliminate the need for `allow_loading_unsigned_plugins`:
1. Create a free organization account on [Grafana.com](https://grafana.com).
2. Generate an Access Policy Token with the `plugins:write` scope.
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

Apache License 2.0. See [LICENSE](LICENSE) for details.
