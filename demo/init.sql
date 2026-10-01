-- Ensure user exists and has monitoring roles
BEGIN
    EXECUTE IMMEDIATE 'CREATE USER grafana_demo IDENTIFIED BY oracle_password';
EXCEPTION
    WHEN OTHERS THEN
        NULL; -- User already created by container startup
END;
/

-- Grant required dictionary access for DBA performance and tablespace dashboards
GRANT CREATE SESSION, RESOURCE TO grafana_demo;
GRANT SELECT ANY DICTIONARY TO grafana_demo;
GRANT SELECT_CATALOG_ROLE TO grafana_demo;
ALTER USER grafana_demo QUOTA UNLIMITED ON USERS;

-- Connect / schema context
ALTER SESSION SET CURRENT_SCHEMA = grafana_demo;

-- Historical metrics table for oracle_analytics.json
CREATE TABLE grafana_demo.metrics_test (
    ts TIMESTAMP DEFAULT SYSTIMESTAMP NOT NULL,
    val NUMBER(10, 2) NOT NULL,
    host VARCHAR2(64) DEFAULT 'oracle-node-1',
    region VARCHAR2(32) DEFAULT 'us-east-1'
);

-- Seed 120 historical data points spanning the last 6 hours
BEGIN
    FOR i IN 1..120 LOOP
        INSERT INTO grafana_demo.metrics_test (ts, val, host, region)
        VALUES (
            SYSTIMESTAMP - NUMTODSINTERVAL(i * 3, 'MINUTE'),
            ROUND(50 + 25 * SIN(i * 0.12) + DBMS_RANDOM.VALUE(-5, 5), 2),
            'oracle-node-' || (MOD(i, 3) + 1),
            'us-east-1'
        );
    END LOOP;
    COMMIT;
END;
/

-- Sample orders table with JSON and CLOB columns (Tier 2 feature showcase)
CREATE TABLE grafana_demo.orders (
    order_id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    customer_name VARCHAR2(100) NOT NULL,
    order_data JSON NOT NULL,
    notes CLOB,
    created_at TIMESTAMP DEFAULT SYSTIMESTAMP
);

INSERT INTO grafana_demo.orders (customer_name, order_data, notes)
VALUES (
    'Enterprise Cloud Logistics',
    JSON('{"order_id": 1001, "total_usd": 15400.00, "status": "COMPLETED", "items": [{"sku": "ORCL-LIC-1", "qty": 4}]}'),
    'Priority enterprise customer order with 24/7 dedicated support SLA.'
);

INSERT INTO grafana_demo.orders (customer_name, order_data, notes)
VALUES (
    'Global Fintech Solutions',
    JSON('{"order_id": 1002, "total_usd": 8250.75, "status": "PROCESSING", "items": [{"sku": "GRAF-ENT-SUB", "qty": 1}]}'),
    'Standard recurring licensing invoice processed via automated billing gateway.'
);

COMMIT;

-- Races table for backward compatibility with legacy examples
CREATE TABLE grafana_demo.races (
    entries INTEGER,
    name VARCHAR2(64) NOT NULL,
    data DATE NOT NULL
);

INSERT ALL
    INTO grafana_demo.races (entries, name, data) VALUES (30, 'Race 1', SYSDATE - 5/24)
    INTO grafana_demo.races (entries, name, data) VALUES (26, 'Race 2', SYSDATE - 4/24)
    INTO grafana_demo.races (entries, name, data) VALUES (30, 'Race 3', SYSDATE - 3/24)
    INTO grafana_demo.races (entries, name, data) VALUES (28, 'Race 4', SYSDATE - 2/24)
    INTO grafana_demo.races (entries, name, data) VALUES (32, 'Race 5', SYSDATE - 1/24)
    INTO grafana_demo.races (entries, name, data) VALUES (34, 'Final', SYSDATE)
SELECT 1 FROM DUAL;

COMMIT;

-- Background synthetic metric generator using DBMS_SCHEDULER
BEGIN
    DBMS_SCHEDULER.CREATE_JOB (
        job_name        => 'GENERATE_GRAFANA_METRICS',
        job_type        => 'PLSQL_BLOCK',
        job_action      => 'BEGIN INSERT INTO grafana_demo.metrics_test (ts, val, host, region) VALUES (SYSTIMESTAMP, ROUND(50 + 25 * SIN(TO_NUMBER(TO_CHAR(SYSTIMESTAMP, ''SS'')) * 0.1) + DBMS_RANDOM.VALUE(-5, 5), 2), ''oracle-node-'' || (MOD(TO_NUMBER(TO_CHAR(SYSTIMESTAMP, ''SS'')), 3) + 1), ''us-east-1''); COMMIT; END;',
        start_date      => SYSTIMESTAMP,
        repeat_interval => 'FREQ=SECONDLY; INTERVAL=10',
        enabled         => TRUE
    );
EXCEPTION
    WHEN OTHERS THEN
        NULL;
END;
/
