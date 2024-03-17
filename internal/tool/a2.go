package tool

import (
	"fmt"
)

// ParseT do short converting on numeric types.
func ParseT[T Parseable](str string) (T, error) {
	var result T
	_, err := fmt.Sscanf(str, "%v", &result)
	return result, err
}

type Parseable interface {
	// NOTE: I didn't check that fmt.Sscanf can accept all these,
	// but it seems like it probably should...
	string | bool | int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr |
		float32 | float64 | complex64 | complex128
}
