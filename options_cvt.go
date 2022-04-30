// Copyright © 2019 Hedzr Yeh.

package cmdr

import "strconv"

func stringSliceToIntSlice(in []string) (out []int) {
	for _, ii := range in {
		if i, err := strconv.Atoi(ii); err == nil {
			out = append(out, i)
		}
	}
	return
}

func stringSliceToInt64Slice(in []string) (out []int64) {
	for _, ii := range in {
		if i, err := strconv.ParseInt(ii, 0, 64); err == nil {
			out = append(out, i)
		}
	}
	return
}

func stringSliceToUint64Slice(in []string) (out []uint64) {
	for _, ii := range in {
		if i, err := strconv.ParseUint(ii, 0, 64); err == nil {
			out = append(out, i)
		}
	}
	return
}

func intSliceToInt64Slice(in []int) (out []int64) {
	for _, ii := range in {
		out = append(out, int64(ii))
	}
	return
}

func int8SliceToInt64Slice(in []int8) (out []int64) {
	for _, ii := range in {
		out = append(out, int64(ii))
	}
	return
}

func int16SliceToInt64Slice(in []int16) (out []int64) {
	for _, ii := range in {
		out = append(out, int64(ii))
	}
	return
}

func int32SliceToInt64Slice(in []int32) (out []int64) {
	for _, ii := range in {
		out = append(out, int64(ii))
	}
	return
}

func intSliceToUint64Slice(in []int) (out []uint64) {
	for _, ii := range in {
		out = append(out, uint64(ii))
	}
	return
}

func int8SliceToUint64Slice(in []int8) (out []uint64) {
	for _, ii := range in {
		out = append(out, uint64(ii))
	}
	return
}

func int16SliceToUint64Slice(in []int16) (out []uint64) {
	for _, ii := range in {
		out = append(out, uint64(ii))
	}
	return
}

func int32SliceToUint64Slice(in []int32) (out []uint64) {
	for _, ii := range in {
		out = append(out, uint64(ii))
	}
	return
}

// func int64SliceToUint64Slice(in []int64) (out []uint64) {
// 	for _, ii := range in {
// 		out = append(out, uint64(ii))
// 	}
// 	return
// }

func int8SliceToIntSlice(in []int8) (out []int) {
	for _, ii := range in {
		out = append(out, int(ii))
	}
	return
}

func int16SliceToIntSlice(in []int16) (out []int) {
	for _, ii := range in {
		out = append(out, int(ii))
	}
	return
}

func int32SliceToIntSlice(in []int32) (out []int) {
	for _, ii := range in {
		out = append(out, int(ii))
	}
	return
}

func int64SliceToIntSlice(in []int64) (out []int) {
	for _, ii := range in {
		out = append(out, int(ii))
	}
	return
}

func int64SliceToUint64Slice(in []int64) (out []uint64) {
	for _, ii := range in {
		out = append(out, uint64(ii))
	}
	return
}

func uintSliceToUint64Slice(in []uint) (out []uint64) {
	for _, ii := range in {
		out = append(out, uint64(ii))
	}
	return
}

// func uint8SliceToUint64Slice(in []uint8) (out []uint64) {
// 	for _, ii := range in {
// 		out = append(out, uint64(ii))
// 	}
// 	return
// }

func uint16SliceToUint64Slice(in []uint16) (out []uint64) {
	for _, ii := range in {
		out = append(out, uint64(ii))
	}
	return
}

func uint32SliceToUint64Slice(in []uint32) (out []uint64) {
	for _, ii := range in {
		out = append(out, uint64(ii))
	}
	return
}

func uintSliceToIntSlice(in []uint) (out []int) {
	for _, ii := range in {
		out = append(out, int(ii))
	}
	return
}

// func uint8SliceToIntSlice(in []uint8) (out []int) {
// 	for _, ii := range in {
// 		out = append(out, int(ii))
// 	}
// 	return
// }

func uint16SliceToIntSlice(in []uint16) (out []int) {
	for _, ii := range in {
		out = append(out, int(ii))
	}
	return
}

func uint32SliceToIntSlice(in []uint32) (out []int) {
	for _, ii := range in {
		out = append(out, int(ii))
	}
	return
}

func uint64SliceToIntSlice(in []uint64) (out []int) {
	for _, ii := range in {
		out = append(out, int(ii))
	}
	return
}

func uintSliceToInt64Slice(in []uint) (out []int64) {
	for _, ii := range in {
		out = append(out, int64(ii))
	}
	return
}

// func uint8SliceToInt64Slice(in []uint8) (out []int64) {
// 	for _, ii := range in {
// 		out = append(out, int64(ii))
// 	}
// 	return
// }

func uint16SliceToInt64Slice(in []uint16) (out []int64) {
	for _, ii := range in {
		out = append(out, int64(ii))
	}
	return
}

func uint32SliceToInt64Slice(in []uint32) (out []int64) {
	for _, ii := range in {
		out = append(out, int64(ii))
	}
	return
}

func uint64SliceToInt64Slice(in []uint64) (out []int64) {
	for _, ii := range in {
		out = append(out, int64(ii))
	}
	return
}
