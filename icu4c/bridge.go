package icu4c

// #cgo pkg-config: icu-i18n icu-uc
// #include "bridge.h"
import "C"

import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/text/language"
)

const bufferSize = 512

var ErrLocalTZ = errors.New("icu4c: Go Local tz name is not supported")

// TZName is the tz database name such as "Asia/Hong_Kong".
// Note that this is different from the name returned by time.Zone().
type TZName string

type DateFormatStyle C.UDateFormatStyle

const (
	DateFormatStyleNone   = C.UDAT_NONE
	DateFormatStyleShort  = C.UDAT_SHORT
	DateFormatStyleMedium = C.UDAT_MEDIUM
	DateFormatStyleLong   = C.UDAT_LONG
	DateFormatStyleFull   = C.UDAT_FULL
)

func FormatDatetime(languageTag language.Tag, tzName TZName, dateStyle DateFormatStyle, timeStyle DateFormatStyle, t time.Time) (out string, err error) {
	if tzName == "Local" {
		err = ErrLocalTZ
		return
	}

	locale := languageTag.String()
	cLocale := C.CString(locale)
	cTZ := C.CString(string(tzName))
	msec := C.double(t.UnixNano() / int64(time.Millisecond))
	resultSize := C.size_t(bufferSize * C.sizeof_char)
	result := (*C.char)(C.malloc(resultSize))

	defer func() {
		C.free(unsafe.Pointer(cLocale))
		C.free(unsafe.Pointer(cTZ))
		C.free(unsafe.Pointer(result))
	}()

	status := C.go_format_datetime(
		cLocale,
		cTZ,
		C.UDateFormatStyle(dateStyle),
		C.UDateFormatStyle(timeStyle),
		msec,
		result,
		resultSize,
	)
	if status != 0 {
		err = fmt.Errorf("icu4c: %v", status)
		return
	}

	out = C.GoString(result)
	return
}
