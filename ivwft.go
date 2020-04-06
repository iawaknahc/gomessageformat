package messageformat

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// IVWFT derives i, v, w, f, t from number according to
// https://unicode.org/reports/tr35/tr35-numbers.html#Operands
func IVWFT(number interface{}) (i, v, w, f, t int, err error) {
	switch value := number.(type) {
	case int8:
		return ivwftInt(absInt(int64(value)))
	case int16:
		return ivwftInt(absInt(int64(value)))
	case int32:
		return ivwftInt(absInt(int64(value)))
	case int64:
		return ivwftInt(absInt(value))
	case int:
		return ivwftInt(absInt(int64(value)))
	case uint8:
		return ivwftInt(uint64(value))
	case uint16:
		return ivwftInt(uint64(value))
	case uint32:
		return ivwftInt(uint64(value))
	case uint64:
		return ivwftInt(value)
	case uint:
		return ivwftInt(uint64(value))
	case float32:
		return ivwftFloat(float64(value))
	case float64:
		return ivwftFloat(value)
	case string:
		return ivwftString(value)
	default:
		err = fmt.Errorf("unsupported type: %T", number)
	}
	return
}

func absInt(number int64) uint64 {
	if number == math.MinInt64 {
		return 1 << 63
	} else if number < 0 {
		return uint64(-number)
	}
	return uint64(number)
}

func parseInt(s string) (i int, err error) {
	if s == "" {
		return
	}
	return strconv.Atoi(s)
}

func ivwftInt(number uint64) (i, v, w, f, t int, err error) {
	return ivwftString(strconv.FormatUint(number, 10))
}

func ivwftFloat(number float64) (i, v, w, f, t int, err error) {
	return ivwftString(strconv.FormatFloat(math.Abs(number), 'f', -1, 64))
}

func ivwftString(number string) (i, v, w, f, t int, err error) {
	if strings.HasPrefix(number, "-") {
		number = number[1:]
	}

	idx := strings.IndexRune(number, '.')
	if idx == -1 {
		i, err = parseInt(number)
		return
	}

	integral := number[0:idx]
	fraction := number[idx+1:]

	i, err = parseInt(integral)
	if err != nil {
		return
	}

	v = len(fraction)
	w = len(strings.TrimRight(fraction, "0"))

	f, err = parseInt(strings.TrimLeft(fraction, "0"))
	if err != nil {
		return
	}

	t, err = parseInt(strings.TrimRight(strings.TrimLeft(fraction, "0"), "0"))
	if err != nil {
		return
	}

	return
}
