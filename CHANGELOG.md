# Oracle Grafana Changelog

## 2.3.0

Feature release introducing out-of-the-box DBA starter dashboards, automatic dashboard discovery in Grafana plugin navigation, and tablespace capacity planning.

### Contributor Sandbox & Evaluation (Tier 4)
- **1-Command Demo Environment (`docker-compose.demo.yml`)**:
  - Full local testbed pairing Oracle Database 23ai Free with Grafana 10.
  - Pre-provisioned datasource (`Oracle-Demo`) and all three bundled dashboards pre-loaded into Grafana.
  - Automated database seed script (`demo/init.sql`) populating historical metrics, sample orders with native JSON and CLOBs, and configuring a `DBMS_SCHEDULER` continuous background metric generator.
  - Added `make demo` and `make demo-down` workflow commands.

### DBA Starter Dashboards (Tier 3)
- **Oracle Database Performance & Health (`oracle_performance.json`)**:
  - Live session tracking: Active user sessions, inactive/idle sessions, and background process count (`v$session`).
  - Cache efficiency: Real-time Buffer Cache Hit Ratio and Library Cache Hit Ratio gauges (`v$sysstat`, `v$librarycache`).
  - Resource bottlenecks: Top 10 non-idle system wait events and full wait class breakdown (`v$system_event`, `v$system_wait_class`).
  - Heavy query analysis: Top 10 SQL statements by elapsed execution time (`v$sqlarea`) with CPU time, buffer gets, and disk reads.
  - Active operations & logs: Long-running operations progress tracking (`v$session_longops`) and hourly redo log switch activity (`v$log_history`).
- **Oracle Tablespace & Storage Capacity (`oracle_tablespaces.json`)**:
  - Global storage KPIs: Total allocated, used, and free storage (GB) with overall utilization gauge.
  - Tablespace metrics: Real-time tablespace space usage percentages (`dba_tablespace_usage_metrics`) and comprehensive allocation breakdown (`dba_tablespaces`, `dba_data_files`, `dba_free_space`).
  - Datafile capacity: Auto-extend status, max expansion limits, and increment sizes (`dba_data_files`).
  - Temporary & Undo health: Temp tablespace allocations (`v$temp_space_header`) and Undo segment status/retention (`dba_undo_extents`).
- **Plugin Dashboard Integration**:
  - Registered all dashboards under `includes` in `plugin.json` (`addToNav: true`) for 1-click importing directly from the Grafana datasource plugin settings.

---

## 2.2.0

Major feature release introducing Oracle Cloud (OCI) Autonomous Database support, Monaco SQL CodeEditor, Oracle Wallet mTLS, dual-column template variables, and rich LOB/JSON data type handling.

### Cloud & Enterprise Connectivity (Tier 2)
- **Oracle Cloud (OCI) & Autonomous Database Support**: Added native TCPS (TLS) encryption and Oracle Wallet configuration for OCI ATP/ADW instances.
- **Oracle Wallet Management**: Added support for auto-login `cwallet.sso` and password-protected `ewallet.p12` PKCS#12 wallet files.
- **LOB & JSON Scanning**: Transparently stream `CLOB`, `NCLOB`, `BLOB`, `XMLTYPE`, and Oracle 21c/23ai/26ai native `JSON` columns into Grafana string fields without memory corruption or truncation.
- **Secret Masking**: Ensured both database and wallet passwords are securely encrypted in `secureJsonData` and masked in logs.

### User Experience & Editor (Tier 1)
- **Monaco SQL CodeEditor**: Replaced plain `<TextArea>` with Grafana's Monaco editor featuring full SQL syntax coloring, line numbers, word-wrap, and auto-indentation.
- **Quick Macros Bar**: Added interactive buttons for 1-click insertion of `$__timeFilter(ts)`, `$__timeFrom()`, `$__timeTo()`, and `$__interval_ms`.
- **Dual-Column Variable Mapping**: Upgraded `metricFindQuery` to support Grafana standard `__text` and `__value` aliases with case-insensitive matching for human-friendly dropdowns.
- **Collapsible Query Preview**: Moved client-interpolated SQL preview into a collapsible drawer to maximize panel editing real estate.

---

## 2.1.2

Feature and stability release introducing native backend macro interpolation for Grafana Alerting and universal Oracle defaults.

### Backend Macro Engine & Alerting
- **Go Backend Macro Engine**: Added native server-side evaluation for `$__timeFilter(column)`, `$__timeFrom()`, `$__timeTo()`, `$__unixEpochFilter(column)`, `$__unixEpochFrom()`, `$__unixEpochTo()`, `$__interval`, and `$__interval_ms`.
- **Full Grafana Alerting Compatibility**: Alert queries execute and evaluate time ranges natively on the server without requiring browser session or frontend JavaScript interpolation.
- **Connection Log Masking**: Fixed empty password replacement in connection debug logging.

### Dashboard & Query Usability
- **Universal Default Query**: Replaced legacy `SYS.races` query with standard `SELECT SYSDATE AS time, 100 AS value FROM DUAL`, preventing `ORA-00942` errors on initial panel creation.
- **Configurable Dashboard Timezone**: Updated `dashboards/oracle_analytics.json` with a `${db_tz}` template variable and native `$__timeFilter(ts_utc)` macros for multi-region compatibility.

---

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
