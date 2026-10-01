package plugin

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

var (
	timeFilterRegex      = regexp.MustCompile(`(?i)\$__timeFilter\(([^)]+)\)`)
	timeFromRegex        = regexp.MustCompile(`(?i)\$__timeFrom\(\)`)
	timeToRegex          = regexp.MustCompile(`(?i)\$__timeTo\(\)`)
	unixEpochFilterRegex = regexp.MustCompile(`(?i)\$__unixEpochFilter\(([^)]+)\)`)
	unixEpochFromRegex   = regexp.MustCompile(`(?i)\$__unixEpochFrom\(\)`)
	unixEpochToRegex     = regexp.MustCompile(`(?i)\$__unixEpochTo\(\)`)
	intervalMsRegex      = regexp.MustCompile(`(?i)\$__interval_ms\b`)
	intervalRegex        = regexp.MustCompile(`(?i)\$__interval\b`)

	// Regex to check if $__from or $__to is inside a numeric division (e.g. $__from / 1000)
	fromDivRegex = regexp.MustCompile(`\$__from(\s*/\s*1000)`)
	toDivRegex   = regexp.MustCompile(`\$__to(\s*/\s*1000)`)
	rawFromRegex = regexp.MustCompile(`\$__from\b`)
	rawToRegex   = regexp.MustCompile(`\$__to\b`)
)

const OracleTimestampFormat = "2006-01-02 15:04:05.000"
const OracleTimestampSQLFormat = "YYYY-MM-DD HH24:MI:SS.FF3"

// FormatOracleTimestamp formats a Go time.Time to an Oracle TO_TIMESTAMP SQL expression in UTC
func FormatOracleTimestamp(t time.Time) string {
	utc := t.UTC()
	return fmt.Sprintf("TO_TIMESTAMP('%s', '%s')", utc.Format(OracleTimestampFormat), OracleTimestampSQLFormat)
}

// InterpolateMacros replaces Grafana query macros with native Oracle SQL expressions.
// This is executed on the Go backend so that Grafana Alerting and background queries
// execute without depending on frontend JavaScript interpolation.
func InterpolateMacros(sql string, timeRange backend.TimeRange, intervalMs int64) string {
	from := timeRange.From
	to := timeRange.To

	// If TimeRange is not populated (e.g. ad-hoc query), default to reasonable fallback
	if from.IsZero() {
		from = time.Now().Add(-6 * time.Hour)
	}
	if to.IsZero() {
		to = time.Now()
	}

	fromTS := FormatOracleTimestamp(from)
	toTS := FormatOracleTimestamp(to)
	fromMilli := from.UnixMilli()
	toMilli := to.UnixMilli()
	fromSec := from.Unix()
	toSec := to.Unix()

	result := sql

	// 1. $__timeFilter(column) -> column >= TO_TIMESTAMP(...) AND column <= TO_TIMESTAMP(...)
	result = timeFilterRegex.ReplaceAllStringFunc(result, func(match string) string {
		submatches := timeFilterRegex.FindStringSubmatch(match)
		if len(submatches) > 1 {
			col := strings.TrimSpace(submatches[1])
			return fmt.Sprintf("%s >= %s AND %s <= %s", col, fromTS, col, toTS)
		}
		return match
	})

	// 2. $__timeFrom() -> TO_TIMESTAMP(...)
	result = timeFromRegex.ReplaceAllString(result, fromTS)

	// 3. $__timeTo() -> TO_TIMESTAMP(...)
	result = timeToRegex.ReplaceAllString(result, toTS)

	// 4. $__unixEpochFilter(column) -> column >= <sec> AND column <= <sec>
	result = unixEpochFilterRegex.ReplaceAllStringFunc(result, func(match string) string {
		submatches := unixEpochFilterRegex.FindStringSubmatch(match)
		if len(submatches) > 1 {
			col := strings.TrimSpace(submatches[1])
			return fmt.Sprintf("%s >= %d AND %s <= %d", col, fromSec, col, toSec)
		}
		return match
	})

	// 5. $__unixEpochFrom() and $__unixEpochTo()
	result = unixEpochFromRegex.ReplaceAllString(result, fmt.Sprintf("%d", fromSec))
	result = unixEpochToRegex.ReplaceAllString(result, fmt.Sprintf("%d", toSec))

	// 6. $__interval_ms and $__interval
	if intervalMs <= 0 {
		intervalMs = 1000
	}
	result = intervalMsRegex.ReplaceAllString(result, fmt.Sprintf("%d", intervalMs))
	result = intervalRegex.ReplaceAllString(result, fmt.Sprintf("%d", intervalMs/1000))

	// 7. Handle $__from and $__to:
	// If pre-processed by frontend into "NUMTODSINTERVAL($__from / 1000, 'SECOND')"
	if fromDivRegex.MatchString(result) {
		result = fromDivRegex.ReplaceAllString(result, fmt.Sprintf("%d$1", fromMilli))
	}
	if toDivRegex.MatchString(result) {
		result = toDivRegex.ReplaceAllString(result, fmt.Sprintf("%d$1", toMilli))
	}

	// Standalone $__from and $__to -> replace directly with Oracle TO_TIMESTAMP
	result = rawFromRegex.ReplaceAllString(result, fromTS)
	result = rawToRegex.ReplaceAllString(result, toTS)

	return result
}
