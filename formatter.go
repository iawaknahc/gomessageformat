package messageformat

import (
	"fmt"
	"strconv"
)

func formatValue(value interface{}) (out string, err error) {
	switch v := value.(type) {
	case int8:
		out = strconv.FormatInt(int64(v), 10)
	case int16:
		out = strconv.FormatInt(int64(v), 10)
	case int32:
		out = strconv.FormatInt(int64(v), 10)
	case int64:
		out = strconv.FormatInt(v, 10)
	case int:
		out = strconv.FormatInt(int64(v), 10)
	case uint8:
		out = strconv.FormatUint(uint64(v), 10)
	case uint16:
		out = strconv.FormatUint(uint64(v), 10)
	case uint32:
		out = strconv.FormatUint(uint64(v), 10)
	case uint64:
		out = strconv.FormatUint(v, 10)
	case uint:
		out = strconv.FormatUint(uint64(v), 10)
	case float32:
		out = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		out = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		out = v
	case bool:
		out = strconv.FormatBool(v)
	default:
		err = fmt.Errorf("unsupported argument type: %T", value)
	}
	return
}

func offsetValue(value interface{}, offset int) (out interface{}, err error) {
	switch v := value.(type) {
	case int8:
		out = int8(int64(v) - int64(offset))
	case int16:
		out = int16(int64(v) - int64(offset))
	case int32:
		out = int32(int64(v) - int64(offset))
	case int64:
		out = int64(int64(v) - int64(offset))
	case int:
		out = int(int64(v) - int64(offset))
	case uint8:
		out = uint8(int64(v) - int64(offset))
	case uint16:
		out = uint16(int64(v) - int64(offset))
	case uint32:
		out = uint32(int64(v) - int64(offset))
	case uint64:
		out = uint64(int64(v) - int64(offset))
	case uint:
		out = uint(int64(v) - int64(offset))
	case float32:
		out = float32(float32(v) - float32(offset))
	case float64:
		out = float64(float64(v) - float64(offset))
	case string:
		var f64 float64
		f64, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return
		}
		out = strconv.FormatFloat(f64-float64(offset), 'f', -1, 64)
	default:
		err = fmt.Errorf("expected numeric type: %T", value)
	}
	return
}

func matchExplicitValue(value interface{}, explicitValue int) (match bool, err error) {
	switch v := value.(type) {
	case int8:
		match = int64(v) == int64(explicitValue)
	case int16:
		match = int64(v) == int64(explicitValue)
	case int32:
		match = int64(v) == int64(explicitValue)
	case int64:
		match = int64(v) == int64(explicitValue)
	case int:
		match = v == explicitValue
	case uint8:
		match = int64(v) == int64(explicitValue)
	case uint16:
		match = int64(v) == int64(explicitValue)
	case uint32:
		match = int64(v) == int64(explicitValue)
	case uint64:
		match = int64(v) == int64(explicitValue)
	case uint:
		match = int64(v) == int64(explicitValue)
	case float32:
		match = float32(v) == float32(explicitValue)
	case float64:
		match = float64(v) == float64(explicitValue)
	case string:
		var f64 float64
		f64, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return
		}
		match = f64 == float64(explicitValue)
	default:
		err = fmt.Errorf("expected numeric type: %T", value)
	}
	return
}
