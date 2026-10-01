package plugin

import (
	"strings"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestInterpolateMacros(t *testing.T) {
	from := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2026, 10, 1, 11, 0, 0, 0, time.UTC)
	tr := backend.TimeRange{From: from, To: to}

	t.Run("timeFilter macro", func(t *testing.T) {
		input := "SELECT * FROM metrics WHERE $__timeFilter(ts)"
		expected := "SELECT * FROM metrics WHERE ts >= TO_TIMESTAMP('2026-10-01 10:00:00.000', 'YYYY-MM-DD HH24:MI:SS.FF3') AND ts <= TO_TIMESTAMP('2026-10-01 11:00:00.000', 'YYYY-MM-DD HH24:MI:SS.FF3')"
		res := InterpolateMacros(input, tr, 1000)
		if res != expected {
			t.Errorf("Expected: %s\nGot: %s", expected, res)
		}
	})

	t.Run("timeFrom and timeTo macros", func(t *testing.T) {
		input := "SELECT * FROM metrics WHERE ts >= $__timeFrom() AND ts <= $__timeTo()"
		expected := "SELECT * FROM metrics WHERE ts >= TO_TIMESTAMP('2026-10-01 10:00:00.000', 'YYYY-MM-DD HH24:MI:SS.FF3') AND ts <= TO_TIMESTAMP('2026-10-01 11:00:00.000', 'YYYY-MM-DD HH24:MI:SS.FF3')"
		res := InterpolateMacros(input, tr, 1000)
		if res != expected {
			t.Errorf("Expected: %s\nGot: %s", expected, res)
		}
	})

	t.Run("standalone raw $__from and $__to", func(t *testing.T) {
		input := "SELECT * FROM metrics WHERE ts BETWEEN $__from AND $__to"
		expected := "SELECT * FROM metrics WHERE ts BETWEEN TO_TIMESTAMP('2026-10-01 10:00:00.000', 'YYYY-MM-DD HH24:MI:SS.FF3') AND TO_TIMESTAMP('2026-10-01 11:00:00.000', 'YYYY-MM-DD HH24:MI:SS.FF3')"
		res := InterpolateMacros(input, tr, 1000)
		if res != expected {
			t.Errorf("Expected: %s\nGot: %s", expected, res)
		}
	})

	t.Run("pre-interpolated frontend NUMTODSINTERVAL", func(t *testing.T) {
		input := "WHERE ts BETWEEN TO_TIMESTAMP('1970-01-01 00:00:00', 'YYYY-MM-DD HH24:MI:SS') + NUMTODSINTERVAL($__from / 1000, 'SECOND') AND TO_TIMESTAMP('1970-01-01 00:00:00', 'YYYY-MM-DD HH24:MI:SS') + NUMTODSINTERVAL($__to / 1000, 'SECOND')"
		res := InterpolateMacros(input, tr, 1000)
		if strings.Contains(res, "$__from") || strings.Contains(res, "$__to") {
			t.Errorf("Expected $__from and $__to to be replaced, got: %s", res)
		}
		if !strings.Contains(res, "1790848800000 / 1000") || !strings.Contains(res, "1790852400000 / 1000") {
			t.Errorf("Expected epoch millisecond replacement, got: %s", res)
		}
	})

	t.Run("unixEpochFilter macro", func(t *testing.T) {
		input := "SELECT * FROM logs WHERE $__unixEpochFilter(created_at)"
		expected := "SELECT * FROM logs WHERE created_at >= 1790848800 AND created_at <= 1790852400"
		res := InterpolateMacros(input, tr, 1000)
		if res != expected {
			t.Errorf("Expected: %s\nGot: %s", expected, res)
		}
	})

	t.Run("interval macros", func(t *testing.T) {
		input := "SELECT $__interval_ms AS interval_ms, $__interval AS interval_sec FROM dual"
		expected := "SELECT 5000 AS interval_ms, 5 AS interval_sec FROM dual"
		res := InterpolateMacros(input, tr, 5000)
		if res != expected {
			t.Errorf("Expected: %s\nGot: %s", expected, res)
		}
	})
}
