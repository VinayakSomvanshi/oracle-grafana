# Oracle Grafana Changelog

## 2.1.1

Maintenance release enabling Grafana Alerting integration.

### Features & Integrations
- **Grafana Alerting Support**: Enabled alerting capability in plugin metadata (`"alerting": true`), allowing Grafana unified alerting rules to evaluate Oracle queries directly in the background.

---

## 2.1.0

Release with expanded signing coverage, automated zip distribution, one-line installer, and sample analytics dashboard.

### Packaging & Distribution
- **Dual Archive Packaging**: Automated both `.tar.gz` and `.zip` distribution packages across all releases.
- **Port 3000 Wildcard Signing**: Expanded cryptographic signature to explicitly cover port 3000 (`http://localhost:3000/`, `http://*:3000/`, `https://*:3000/`) allowing out-of-the-box loading on standard Grafana instances without `allow_loading_unsigned_plugins`.
- **One-Command CLI Installation**: Fully supported `grafana cli plugins install --pluginUrl ...` directly from GitHub releases.
- **Automated Installer**: Added standalone `install.sh` with automatic Grafana environment detection, permission configuration, and restart instructions.

### Documentation & Dashboards
- **Sample Analytics Dashboard**: Added a comprehensive sample dashboard (`dashboards/oracle_analytics.json`) featuring KPI stat cards, analytic moving average time series, status gauges, raw tables with health badges, and live Oracle data dictionary monitoring.
- **Uninstallation Guide**: Added complete uninstallation procedures for Grafana CLI, manual directories, and Docker containers.
- **ASCII & Quality Standards**: Cleaned all documentation to 100% clean ASCII.

---

## 2.0.0

Modernized release with strict read-only security engine.

### Security & Compliance
- **3-Tier Read-Only Defense Engine**:
  - Implemented SQL Lexer/Sanitizer that strips comments and string literals, and enforces statements starting strictly with SELECT or WITH.
  - Prohibited statement chaining and multi-statement execution via internal semicolons.
  - Prohibited mutating DDL/DML statements (INSERT, UPDATE, DELETE, DROP, ALTER, TRUNCATE, MERGE, CREATE, RENAME, GRANT, REVOKE).
  - Prohibited PL/SQL execution (BEGIN, DECLARE, EXECUTE, CALL).
  - Pre-flight ALTER SESSION SET TRANSACTION READ ONLY on each query execution to enforce database kernel-level mutation blocks (ORA-01456).
- **Context Cancellation**: Fully wired context.Context to PrepareContext and QueryContext to abort queries in Oracle when Grafana panels timeout or tabs close.
- **Circuit Breaker**: Added a 50,000 maximum row limit with frame warning notices (data.Notice) to prevent Out-Of-Memory (OOM) crashes.
- **Vulnerability Remediation**: Updated dependencies to resolve all known CVEs (govulncheck verified: 0 CVEs).

### Performance & Engine
- **Native Type Scanner**: Replaced stringified column scanning with native typed DataFrame builder (*float64, *time.Time, *string, *bool).
- **Modern SDK**: Upgraded github.com/grafana/grafana-plugin-sdk-go to v0.296.5 and go-ora to v2.9.0.
- **Pure-Go Compilation**: Completely CGO-free (CGO_ENABLED=0) cross-compilation for Linux, macOS, and Windows.
- **Connection Pooling**: Configured max open, idle, and lifetime connection pool settings.

### Operations
- Added production Oracle setup script (oracle/setup_production.sql) with resource profile limits (c_grafana_profile).
- Added automated CI/CD pipeline (.github/workflows/ci.yml) with automated vulnerability checks, cross-compilation, and Grafana signing.
- Added production host installer (install_production.sh).

---

## 1.0.0 (Original Upstream Release)

Initial release by Open source community as a Datasource with internal backend support.
