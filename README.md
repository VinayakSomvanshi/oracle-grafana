# Oracle Database Data Source for Grafana

[![CI & Release Pipeline](https://github.com/VinayakSomvanshi/oracle-grafana/actions/workflows/ci.yml/badge.svg)](https://github.com/VinayakSomvanshi/oracle-grafana/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

An enterprise-ready, zero-configuration, strictly read-only Oracle Database data source plugin for Grafana OSS and Enterprise (versions 10.x, 11.x, and newer).

This plugin connects Grafana directly to any Oracle Database instance (11g, 12c, 18c, 19c, 21c, 23ai, and 26ai) with zero database-side modifications or specialized DBA privileges.

---

## Key Capabilities

- Zero Database Administration: Connects with standard database credentials. Does not require DBA rights, custom schemas, stored procedures, or database-side scripts.
- Pure Go Architecture: Built on the pure-Go driver (`sijms/go-ora/v2`) with `CGO_ENABLED=0`. Operates without Oracle Instant Client packages, C runtime libraries, or native system dependencies.
- Autonomous Read-Only Protection: Multi-tier safety engine guarantees zero data modification risk:
  - In-Plugin Query Lexer: Strips comments and asserts queries begin strictly with `SELECT` or `WITH`.
  - Kernel-Level Session Isolation: Executes `ALTER SESSION SET TRANSACTION READ ONLY` on checked-out database connections. Any attempted data modification (`INSERT`, `UPDATE`, `DELETE`, `MERGE`, `DROP`, `ALTER`, PL/SQL blocks) is rejected directly by the Oracle database engine with `ORA-01456`.
  - Transaction Rollback and Connection Sanitization: Issues an explicit `ROLLBACK` before returning connections to the pool, preventing session pollution.
- High-Performance Native Typing: Directly maps Oracle data types (`NUMBER`, `FLOAT`, `DATE`, `TIMESTAMP`, `VARCHAR2`, `CHAR`) into Grafana DataFrame vectors for instant rendering.
- Circuit Breaker: Enforces a 50,000-row safety limit with Grafana DataFrame warning notices to prevent host Out-Of-Memory (OOM) conditions.
- Context Cancellation: Automatically maps dashboard tab closures and panel timeouts to Oracle query cancellation signals.
- Multi-Platform Support: Pre-compiled binaries for Linux (amd64, arm64), macOS (Intel, Apple Silicon), and Windows (amd64).

---

## Installation

### Option 1: Direct Pre-Compiled Release Download (Recommended)

Pre-built release archives contain all cross-platform backends and the frontend distribution.

1. Download the latest release archive:
   ```bash
   # Download tarball release
   curl -fSL -o oracle-grafana-datasource-2.0.0.tar.gz https://github.com/VinayakSomvanshi/oracle-grafana/releases/latest/download/oracle-grafana-datasource-2.0.0.tar.gz
   ```

2. Extract into your Grafana plugins directory (default: `/var/lib/grafana/plugins`):
   ```bash
   sudo mkdir -p /var/lib/grafana/plugins/oracle-grafana-datasource
   sudo tar -xzf oracle-grafana-datasource-2.0.0.tar.gz -C /var/lib/grafana/plugins/oracle-grafana-datasource --strip-components=1
   sudo chown -R grafana:grafana /var/lib/grafana/plugins/oracle-grafana-datasource
   ```

3. Restart the Grafana server:
   ```bash
   # Systemd service
   sudo systemctl restart grafana-server

   # Docker container
   docker restart grafana
   ```

### Option 2: Automated Local Installation Script

If running on the local machine where the repository was cloned:
```bash
sudo ./install_production.sh
```
This script detects the Grafana plugin directory, builds or syncs release artifacts, configures file permissions, and prompts to restart the service.

### Option 3: Building from Source

Prerequisites:
- Go 1.22 or higher
- Node.js 20 or higher
- Yarn or NPM
- Make

```bash
git clone https://github.com/VinayakSomvanshi/oracle-grafana.git
cd oracle-grafana
make dist
```

---

## Plugin Signature and Verification

Official release archives are cryptographically signed using Grafana Labs' official signing tool (`@grafana/sign-plugin`) configured with universal root URL patterns (`http://*`, `https://*`). This signature allows the plugin to load on any hostname or IP address without signature warning barriers.

If compiling custom unsigned builds, enable unsigned plugin loading in `/etc/grafana/grafana.ini`:
```ini
[plugins]
allow_loading_unsigned_plugins = oracle-grafana-datasource
```
Then restart `grafana-server`.

---

## Configuration in Grafana

1. Open Grafana in your web browser.
2. Navigate to **Connections > Data Sources > Add data source**.
3. Search for **Oracle** and select the plugin.
4. Fill in the connection settings:
   - **Host**: Database server IP or hostname (e.g., `10.0.0.25` or `db.internal.net`).
   - **Port**: Oracle listener port (default: `1521`).
   - **Service / SID**: Oracle service name (e.g., `ORCL`, `FREEPDB1`) or SID.
   - **User**: Any standard Oracle database user.
   - **Password**: User password.
   - **Connection Pool Settings**:
     - **Max Open Connections**: Maximum active connections in pool (default: `5`).
     - **Max Idle Connections**: Maximum idle connections retained (default: `2`).
     - **Connection Max Lifetime**: Maximum connection lifetime in seconds (default: `14400` / 4 hours).
5. Click **Save & Test**. Grafana performs an immediate health check (`SELECT 1 FROM DUAL` within a read-only session) to verify listener reachability, authentication, and read-only policy enforcement.

---

## Querying Oracle in Grafana Panels

The plugin supports both Time Series and Table visualizations, as well as dynamic Dashboard Template Variables.

### 1. Time Series Visualizations

For Time Series panels, the query must return a timestamp column and at least one numeric metric column. Optionally, a string column can be included to represent the metric name or series name.

#### Example: Active Database Sessions by Status
```sql
SELECT
  SYSTIMESTAMP AS "time",
  status AS "metric",
  COUNT(*) AS "value"
FROM v$session
WHERE type = 'USER'
GROUP BY status, SYSTIMESTAMP
ORDER BY "time" ASC
```

#### Example: Tablespace Usage with Grafana Time Filters
```sql
SELECT
  collection_time AS "time",
  tablespace_name AS "metric",
  used_space_pct AS "value"
FROM dba_tablespace_usage_metrics
WHERE collection_time >= TO_TIMESTAMP('${__from:date:YYYY-MM-DD HH24:MI:SS}', 'YYYY-MM-DD HH24:MI:SS')
  AND collection_time <= TO_TIMESTAMP('${__to:date:YYYY-MM-DD HH24:MI:SS}', 'YYYY-MM-DD HH24:MI:SS')
ORDER BY collection_time ASC
```

#### Supported Time Range Macros
- `${__from:date:YYYY-MM-DD HH24:MI:SS}`: Start of the Grafana dashboard time picker range formatted for Oracle `TO_TIMESTAMP`.
- `${__to:date:YYYY-MM-DD HH24:MI:SS}`: End of the Grafana dashboard time picker range formatted for Oracle `TO_TIMESTAMP`.
- `$__from`: Start time in Unix epoch milliseconds.
- `$__to`: End time in Unix epoch milliseconds.

### 2. Table Visualizations

For Table panels, write standard SQL queries. The plugin formats all returned rows into typed columns.

#### Example: Tablespace Allocation and Capacity
```sql
SELECT
  tablespace_name,
  ROUND(SUM(bytes) / (1024 * 1024), 2) AS total_mb,
  ROUND(SUM(maxbytes) / (1024 * 1024), 2) AS max_mb,
  autoextensible
FROM dba_data_files
GROUP BY tablespace_name, autoextensible
ORDER BY tablespace_name ASC
```

#### Example: Longest Running Active Queries
```sql
SELECT
  sid,
  serial#,
  username,
  sql_id,
  last_call_et AS elapsed_seconds,
  status
FROM v$session
WHERE status = 'ACTIVE'
  AND type = 'USER'
ORDER BY last_call_et DESC
```

### 3. Dashboard Template Variables

Use Oracle queries to dynamically populate dashboard drop-down filters.

1. Navigate to **Dashboard Settings > Variables > New variable**.
2. Set Variable Type to **Query** and select your Oracle data source.
3. Enter the query in the **Query** field:

#### Example: Tablespace Dropdown Filter
```sql
SELECT DISTINCT tablespace_name FROM dba_tablespaces ORDER BY tablespace_name ASC
```

#### Example: Database Schema/User Filter
```sql
SELECT username FROM all_users ORDER BY username ASC
```

You can then reference the variable in panel queries using `$variable_name` or `${variable_name}`:
```sql
SELECT
  table_name,
  num_rows,
  last_analyzed
FROM all_tables
WHERE owner = '${schema_name}'
ORDER BY table_name ASC
```

---

## Safety and Read-Only Enforcement

The plugin autonomously prevents any data modification on the database:

1. Pre-Flight Lexer Inspection:
   Incoming queries are sanitized by removing whitespace and comments. Queries that do not begin with `SELECT` or `WITH` are rejected at the application layer with a descriptive validation error.

2. Session-Level Read-Only Transaction:
   Before query execution, the plugin obtains an isolated connection from the pool and executes:
   ```sql
   ALTER SESSION SET TRANSACTION READ ONLY
   ```
   Even if the configured user has `DBA`, `DROP`, or `DELETE` privileges, the Oracle database kernel prevents any data modification with:
   ```text
   ORA-01456: may not perform insert/delete/update operation inside a READ ONLY transaction
   ```

3. Session State Rollback:
   When query execution concludes, an explicit `ROLLBACK` is issued to clear transaction state before the connection returns to the connection pool.

---

## Compatibility Matrix

### Oracle Database Versions
- Oracle Database 11g Release 2
- Oracle Database 12c (Release 1 and 2, CDB and non-CDB)
- Oracle Database 18c
- Oracle Database 19c (Enterprise, Standard, XE)
- Oracle Database 21c (Innovation Release)
- Oracle Database 23ai (Free, Enterprise, Cloud)
- Oracle Database 26ai

### Grafana Platforms
- Grafana OSS 10.x, 11.x, and newer
- Grafana Enterprise 10.x, 11.x, and newer
- Grafana Cloud (via Private Data Source Connect or self-hosted relay)

### Operating Systems & Architectures
- Linux: `x86_64` (`amd64`), `aarch64` (`arm64`)
- macOS: Intel (`amd64`), Apple Silicon (`arm64`)
- Windows: `x86_64` (`amd64`)

---

## Development and Testing

```bash
# Run backend unit tests
make test

# Cross-compile all backends
make build-all-backends

# Run Go security vulnerability scanner
govulncheck ./...

# Package complete release archive
make dist

# Clean build artifacts
make clean
```

---

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE) for complete licensing information.
