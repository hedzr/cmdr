package tool

import (
	"strings"
)

// ToBool translate a value (int, bool, string) to boolean
func ToBool(val any, defaultVal ...bool) (ret bool) {
	if val != nil {
		if v, ok := val.(bool); ok {
			return v
		}
		switch v := val.(type) {
		case int:
			return v != 0
		case int8:
			return v != 0
		case int16:
			return v != 0
		case int32:
			return v != 0
		case int64:
			return v != 0
		case uint:
			return v != 0
		case uint8:
			return v != 0
		case uint16:
			return v != 0
		case uint32:
			return v != 0
		case uint64:
			return v != 0
		}
		if v, ok := val.(string); ok {
			return toBoolClassical(v, defaultVal...)
		}
	}
	return toBoolDefVal(defaultVal...)
}

// func isZero[T slog.Integers | slog.Uintegers](v T) bool {
// 	return v == 0
// }

func toBoolDefVal(defaultVal ...bool) (ret bool) {
	for _, vv := range defaultVal {
		ret = vv
	}
	return
}

func StringToBool(val string, defaultVal ...bool) (ret bool) {
	return toBoolClassical(val, defaultVal...)
}

func toBoolClassical(val string, defaultVal ...bool) (ret bool) {
	// ret = ToBool(val, defaultVal...)
	switch strings.ToLower(val) {
	case "1", "y", "t", "yes", "true", "ok", "on":
		ret = true
	case "":
		return toBoolDefVal(defaultVal...)
	}
	return
}
