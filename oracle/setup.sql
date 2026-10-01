-- Run as SYSTEM connected to FREEPDB1: sqlplus system@localhost:1521/FREEPDB1
-- CHANGE THE TEST PASSWORD.
CREATE USER grafana IDENTIFIED BY "YourSecurePasswordHere123!" QUOTA 50M ON USERS;
GRANT CREATE SESSION, CREATE TABLE TO grafana;
CREATE TABLE grafana.metrics_test (ts DATE, val NUMBER);
INSERT INTO grafana.metrics_test
  SELECT SYSDATE - LEVEL/1440, ROUND(DBMS_RANDOM.VALUE(1,100)) FROM dual CONNECT BY LEVEL <= 120;
COMMIT;

-- Optional: continuous test data (one row per second) and a purge job
BEGIN
  DBMS_SCHEDULER.CREATE_JOB(job_name => 'INSERT_TEST_ROW', job_type => 'PLSQL_BLOCK',
    job_action => 'BEGIN INSERT INTO grafana.metrics_test VALUES (SYSDATE, ROUND(DBMS_RANDOM.VALUE(1,100))); COMMIT; END;',
    repeat_interval => 'FREQ=SECONDLY;INTERVAL=1', enabled => TRUE);
  DBMS_SCHEDULER.SET_ATTRIBUTE('INSERT_TEST_ROW', 'logging_level', DBMS_SCHEDULER.LOGGING_OFF);
END;
/
BEGIN
  DBMS_SCHEDULER.CREATE_JOB(job_name => 'PURGE_TEST_ROWS', job_type => 'PLSQL_BLOCK',
    job_action => 'BEGIN DELETE FROM grafana.metrics_test WHERE ts < SYSDATE - 1; COMMIT; END;',
    repeat_interval => 'FREQ=HOURLY', enabled => TRUE);
END;
/
