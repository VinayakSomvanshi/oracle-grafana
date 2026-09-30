-- ============================================================================
-- PRODUCTION READ-ONLY DATABASE SETUP FOR GRAFANA
-- Target: Oracle Database 19c / 21c / 23ai / 26ai (CDB/PDB)
-- Run as SYSDBA or SYSTEM connected to the application PDB (e.g. FREEPDB1)
-- ============================================================================

-- 1. Create a Resource Profile (Prevents Runaway Queries & Resource Exhaustion)
BEGIN
  EXECUTE IMMEDIATE 'DROP PROFILE c_grafana_profile CASCADE';
EXCEPTION
  WHEN OTHERS THEN NULL;
END;
/

CREATE PROFILE c_grafana_profile LIMIT
  SESSIONS_PER_USER          15       -- Max concurrent connections from Grafana pool
  CPU_PER_CALL               3000     -- 30 seconds CPU time max per call (in 100ths of sec)
  LOGICAL_READS_PER_CALL     250000   -- Limit memory/buffer cache churn
  IDLE_TIME                  15       -- Terminate idle sessions after 15 minutes
  CONNECT_TIME               120      -- Max session duration in minutes (reconnects cleanly)
  FAILED_LOGIN_ATTEMPTS      5        -- Lock account after 5 failed tries
  PASSWORD_LIFE_TIME         90;      -- Force key rotation every 90 days

-- 2. Create the Dedicated Read-Only Service Account
-- IMPORTANT: Change the password below to a secure, randomly generated 32+ char secret.
CREATE USER grafana_svc IDENTIFIED BY "ChangeMe_To_A_Secure_Password_123!"
  DEFAULT TABLESPACE USERS
  TEMPORARY TABLESPACE TEMP
  PROFILE c_grafana_profile
  QUOTA 0 ON USERS;

-- 3. Grant Minimum Necessary Privileges (ZERO DDL, ZERO DML)
GRANT CREATE SESSION TO grafana_svc;

-- Optional: Allow inspecting data dictionary if you want schema browsing
GRANT SELECT_CATALOG_ROLE TO grafana_svc;

-- 4. Explicit Object Grants (Principle of Least Privilege)
-- Grant SELECT ONLY on the specific tables or views required for Grafana dashboards:
-- Example:
-- GRANT SELECT ON app_schema.orders TO grafana_svc;
-- GRANT SELECT ON app_schema.metrics TO grafana_svc;

-- 5. Compliance: Oracle Unified Auditing Policy (Optional, for Enterprise Auditing)
-- Track every query executed by the Grafana service account
/*
CREATE AUDIT POLICY grafana_query_audit
  ACTIONS ALL
  WHEN 'SYS_CONTEXT(''USERENV'', ''SESSION_USER'') = ''GRAFANA_SVC'''
  EVALUATE PER STATEMENT;

AUDIT POLICY grafana_query_audit;
*/

PROMPT Production Grafana Read-Only setup completed successfully.
