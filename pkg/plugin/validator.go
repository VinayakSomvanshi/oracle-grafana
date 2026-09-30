package plugin

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	blockCommentRegex = regexp.MustCompile(`/\*[\s\S]*?\*/`)
	lineCommentRegex  = regexp.MustCompile(`--[^\r\n]*`)
	stringLitRegex    = regexp.MustCompile(`'([^']|'')*'`)
	forbiddenWordRegex = regexp.MustCompile(`(?i)\b(INSERT|UPDATE|DELETE|MERGE|DROP|ALTER|TRUNCATE|CREATE|RENAME|GRANT|REVOKE|BEGIN|DECLARE|EXEC|EXECUTE|CALL|LOCK)\b`)
	leadingWordRegex   = regexp.MustCompile(`(?i)^[ \t\r\n]*(SELECT|WITH)\b`)
)

// StripSQLCommentsAndLiterals removes comments and string literals so that
// keyword analysis only operates on actual SQL code.
func StripSQLCommentsAndLiterals(sql string) string {
	// Strip block comments
	noBlock := blockCommentRegex.ReplaceAllString(sql, " ")
	// Strip line comments
	noLine := lineCommentRegex.ReplaceAllString(noBlock, " ")
	// Strip string literals
	noLits := stringLitRegex.ReplaceAllString(noLine, "''")
	return strings.TrimSpace(noLits)
}

// ValidateReadOnlyQuery ensures that the given SQL query is strictly read-only,
// contains only a single statement, starts with SELECT or WITH, and does not contain
// any mutating DDL, DML, or PL/SQL operations.
func ValidateReadOnlyQuery(sql string) error {
	sanitized := StripSQLCommentsAndLiterals(sql)
	if sanitized == "" {
		return fmt.Errorf("query is empty or contains only comments")
	}

	// Remove one optional trailing semicolon
	trimmed := strings.TrimRight(sanitized, "; \t\r\n")

	// Disallow semicolons inside the query to prevent statement chaining / multi-query attacks
	if strings.Contains(trimmed, ";") {
		return fmt.Errorf("multiple SQL statements or embedded semicolons are not permitted")
	}

	// Must start with SELECT or WITH
	if !leadingWordRegex.MatchString(trimmed) {
		return fmt.Errorf("only SELECT and WITH statements are permitted")
	}

	// Check for forbidden keywords
	if match := forbiddenWordRegex.FindString(trimmed); match != "" {
		return fmt.Errorf("forbidden SQL statement or operation detected: %s", strings.ToUpper(match))
	}

	return nil
}
