# Oracle Grafana Changelog

## 2.0.0

Modernized release maintained by Vinayak Somvanshi.

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
