package text

import (
	"strings"
)

// ToBool translate a value (int, bool, string) to boolean
func ToBool(val any, defaultVal ...bool) (ret bool) {
	if val != nil {
		if v, ok := val.(bool); ok {
			return v
		}
		switch val.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:
			return val != 0
		}
		if v, ok := val.(string); ok {
			return toBool(v, defaultVal...)
		}
	}
	for _, vv := range defaultVal {
		ret = vv
	}
	return
}

func StringToBool(val string, defaultVal ...bool) (ret bool) {
	return toBool(val, defaultVal...)
}

func toBool(val string, defaultVal ...bool) (ret bool) {
	// ret = ToBool(val, defaultVal...)
	switch strings.ToLower(val) {
	case "1", "y", "t", "yes", "true", "ok", "on":
		ret = true
	case "":
		for _, vv := range defaultVal {
			ret = vv
		}
	}
	return
}
