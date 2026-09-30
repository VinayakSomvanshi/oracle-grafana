package plugin

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type ColumnKind int

const (
	KindUnknown ColumnKind = iota
	KindNumber
	KindTime
	KindString
	KindBool
)

type TypedColumnBuilder struct {
	Name string
	Kind ColumnKind

	// Slices for each kind
	Numbers []*float64
	Times   []*time.Time
	Strings []*string
	Bools   []*bool
}

func NewColumnBuilder(col *sql.ColumnType) *TypedColumnBuilder {
	name := col.Name()
	dbType := strings.ToUpper(col.DatabaseTypeName())

	kind := KindString
	if strings.Contains(dbType, "INT") ||
		strings.Contains(dbType, "NUM") ||
		strings.Contains(dbType, "FLOAT") ||
		strings.Contains(dbType, "DOUBLE") ||
		strings.Contains(dbType, "DECIMAL") ||
		strings.Contains(dbType, "REAL") {
		kind = KindNumber
	} else if strings.Contains(dbType, "DATE") || strings.Contains(dbType, "TIME") {
		kind = KindTime
	} else if strings.Contains(dbType, "BOOL") {
		kind = KindBool
	}

	return &TypedColumnBuilder{
		Name: name,
		Kind: kind,
	}
}

func (b *TypedColumnBuilder) Append(raw interface{}) {
	if raw == nil {
		switch b.Kind {
		case KindNumber:
			b.Numbers = append(b.Numbers, nil)
		case KindTime:
			b.Times = append(b.Times, nil)
		case KindBool:
			b.Bools = append(b.Bools, nil)
		default:
			b.Strings = append(b.Strings, nil)
		}
		return
	}

	switch b.Kind {
	case KindNumber:
		var num float64
		switch v := raw.(type) {
		case float64:
			num = v
		case float32:
			num = float64(v)
		case int64:
			num = float64(v)
		case int32:
			num = float64(v)
		case int:
			num = float64(v)
		case []byte:
			parsed, err := strconv.ParseFloat(string(v), 64)
			if err == nil {
				num = parsed
			}
		case string:
			parsed, err := strconv.ParseFloat(v, 64)
			if err == nil {
				num = parsed
			}
		default:
			str := fmt.Sprintf("%v", v)
			parsed, err := strconv.ParseFloat(str, 64)
			if err == nil {
				num = parsed
			}
		}
		b.Numbers = append(b.Numbers, &num)

	case KindTime:
		var t time.Time
		switch v := raw.(type) {
		case time.Time:
			t = v
		case string:
			parsed, err := parseTimeString(v)
			if err == nil {
				t = parsed
			}
		case []byte:
			parsed, err := parseTimeString(string(v))
			if err == nil {
				t = parsed
			}
		default:
			str := fmt.Sprintf("%v", v)
			parsed, err := parseTimeString(str)
			if err == nil {
				t = parsed
			}
		}
		b.Times = append(b.Times, &t)

	case KindBool:
		var bVal bool
		switch v := raw.(type) {
		case bool:
			bVal = v
		case int64:
			bVal = (v != 0)
		case string:
			bVal = (strings.EqualFold(v, "true") || v == "1" || strings.EqualFold(v, "y"))
		default:
			bVal = false
		}
		b.Bools = append(b.Bools, &bVal)

	default: // KindString
		var str string
		switch v := raw.(type) {
		case []byte:
			str = string(v)
		case string:
			str = v
		case time.Time:
			str = v.Format(time.RFC3339)
		default:
			str = fmt.Sprintf("%v", v)
		}
		b.Strings = append(b.Strings, &str)
	}
}

func parseTimeString(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02-JAN-06 03.04.05.000000 PM",
		"02-JAN-06",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse date: %s", s)
}

func (b *TypedColumnBuilder) ToField() *data.Field {
	switch b.Kind {
	case KindNumber:
		return data.NewField(b.Name, nil, b.Numbers)
	case KindTime:
		return data.NewField(b.Name, nil, b.Times)
	case KindBool:
		return data.NewField(b.Name, nil, b.Bools)
	default:
		return data.NewField(b.Name, nil, b.Strings)
	}
}
