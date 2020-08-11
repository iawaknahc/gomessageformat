package messageformat

import (
	"fmt"
	"strconv"
)

type argumentMinusOffset struct {
	Name  string
	Value interface{}
}

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
