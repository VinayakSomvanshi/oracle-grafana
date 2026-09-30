# Security & Compliance Architecture

This fork of the Oracle Grafana DataSource has been re-engineered for **enterprise security**, **zero-trust database access**, and **strict read-only compliance**.

---

## 1. Multi-Tier Defense-in-Depth Model

To guarantee that Grafana can **never harm, mutate, alter, or drop data**, the plugin implements a 3-layer security defense:

```
[ User Panel Query in Grafana ]
               │
               ▼
┌────────────────────────────────────────────────────────┐
│ Layer 1: Application-Level AST / Lexer Validation      │
│  - Strips comments (--, /* */) and string literals     │
│  - Enforces statement begins with SELECT or WITH       │
│  - Forbids multi-query chaining (semicolon injection)  │
│  - Blocks DDL/DML tokens (INSERT, UPDATE, DELETE, ...) │
│  - Blocks PL/SQL blocks (BEGIN, DECLARE, EXECUTE)      │
└────────────────────────────────────────────────────────┘
               │ (Pass)
               ▼
┌────────────────────────────────────────────────────────┐
│ Layer 2: Oracle Database Engine Session Lock           │
│  - Issues: ALTER SESSION SET TRANSACTION READ ONLY     │
│  - Oracle kernel rejects any state change with:        │
│    ORA-01456: may not perform insert/delete/update     │
└────────────────────────────────────────────────────────┘
               │
               ▼
┌────────────────────────────────────────────────────────┐
│ Layer 3: Database Role & Object Permissions            │
│  - Service user `grafana_svc` has only CREATE SESSION  │
│  - ZERO DDL privileges (no CREATE TABLE, DROP, ALTER)  │
│  - Explicit SELECT grants ONLY on approved objects     │
│  - Resource Profile limits CPU and Logical IO churn    │
└────────────────────────────────────────────────────────┘
```

---

## 2. Resource Management & Denial-of-Service Protection

* **Context Cancellation**: All SQL queries use `PrepareContext(ctx)` and `QueryContext(ctx)`. When a dashboard panel times out or a browser tab closes, the query in Oracle is aborted immediately.
* **Row-Limit Circuit Breaker**: Hard-capped at 50,000 rows. Queries exceeding this limit are truncated, memory is protected, and a visual warning is attached to the panel without crashing Grafana.
* **Connection Pooling**: Managed via Go's `sql.DB`:
  * `MaxOpenConns`: 25
  * `MaxIdleConns`: 5
  * `ConnMaxLifetime`: 30 minutes
  * `ConnMaxIdleTime`: 5 minutes

---

## 3. Dependency Security & Scanning

* Pure Go driver (`github.com/sijms/go-ora/v2`) — zero Oracle Instant Client or CGO dependencies.
* Continuous vulnerability scanning via `govulncheck` in the CI pipeline.
* Zero known CVEs in current dependency graph.

---

## 4. Cryptographic Signing

The plugin is designed to be cryptographically signed using Grafana's official signing service:
```bash
export GRAFANA_ACCESS_POLICY_TOKEN="<your-token>"
make sign
```
Signed plugins generate a `MANIFEST.txt` and do not require `allow_loading_unsigned_plugins`.
