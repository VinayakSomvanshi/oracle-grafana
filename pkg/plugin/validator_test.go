package plugin

import "testing"

func TestValidateReadOnlyQuery(t *testing.T) {
	validQueries := []string{
		"SELECT * FROM dual",
		"select id, name from users where age > 10",
		"WITH summary AS (SELECT dept, count(*) as cnt FROM emp GROUP BY dept) SELECT * FROM summary",
		"-- leading comment\nSELECT 1 FROM dual",
		"/* multi-line\ncomment */ SELECT sysdate FROM dual;",
		"SELECT col_update_date FROM audit_log",
	}

	for _, q := range validQueries {
		if err := ValidateReadOnlyQuery(q); err != nil {
			t.Errorf("expected query to be valid, got error: %v for query: %s", err, q)
		}
	}

	invalidQueries := []string{
		"INSERT INTO users VALUES (1, 'alice')",
		"UPDATE users SET name='bob' WHERE id=1",
		"DELETE FROM users",
		"DROP TABLE users",
		"ALTER TABLE users ADD col1 VARCHAR2(10)",
		"TRUNCATE TABLE users",
		"CREATE TABLE hack (id INT)",
		"SELECT 1 FROM dual; DROP TABLE users",
		"BEGIN DBMS_LOCK.SLEEP(10); END;",
		"DECLARE x number; BEGIN x := 1; END;",
		"EXECUTE my_proc()",
		"GRANT DBA TO grafana",
		"REVOKE ALL ON users FROM grafana",
		"",
		"   ",
		"-- just a comment",
	}

	for _, q := range invalidQueries {
		if err := ValidateReadOnlyQuery(q); err == nil {
			t.Errorf("expected query to be rejected as invalid: %s", q)
		}
	}
}

func TestLiteralWithForbiddenWord(t *testing.T) {
	q := "SELECT 'INSERT INTO dummy' AS col1, 'DELETE' AS col2 FROM dual"
	if err := ValidateReadOnlyQuery(q); err != nil {
		t.Fatalf("expected literal query to be valid, got: %v", err)
	}
}
