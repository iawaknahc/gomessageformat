package icu4c

import (
	"testing"
	"time"

	"golang.org/x/text/language"
)

func TestFormatDate(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	tz := TZName("Asia/Hong_Kong")
	languageTag := language.Make("zh-Hant-HK")

	expected := "2009年11月11日星期三"

	actual, err := FormatDatetime(languageTag, tz, DateFormatStyleFull, DateFormatStyleNone, now)
	if err != nil {
		t.Errorf("err: %v", err)
	} else if actual != expected {
		t.Errorf("%v != %v", actual, expected)
	}
}

func TestFormatTime(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	tz := TZName("Asia/Hong_Kong")
	languageTag := language.Make("zh-Hant-HK")

	expected := "上午7:00:00 [香港標準時間]"

	actual, err := FormatDatetime(languageTag, tz, DateFormatStyleNone, DateFormatStyleFull, now)
	if err != nil {
		t.Errorf("err: %v", err)
	} else if actual != expected {
		t.Errorf("%v != %v", actual, expected)
	}
}

func TestFormatDatetime(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	tz := TZName("Asia/Hong_Kong")
	languageTag := language.Make("zh-Hant-HK")

	expected := "2009年11月11日星期三 上午7:00:00 [香港標準時間]"

	actual, err := FormatDatetime(languageTag, tz, DateFormatStyleFull, DateFormatStyleFull, now)
	if err != nil {
		t.Errorf("err: %v", err)
	} else if actual != expected {
		t.Errorf("%v != %v", actual, expected)
	}
}

func TestFormatDatetimeLocal(t *testing.T) {
	now := time.Now()
	tz := now.Location().String()
	languageTag := language.Make("zh-Hant-HK")

	_, err := FormatDatetime(languageTag, TZName(tz), DateFormatStyleFull, DateFormatStyleFull, now)
	if err == nil {
		t.Errorf("expected error")
	} else if err.Error() != "icu4c: Go Local tz name is not supported" {
		t.Errorf("unexpected error: %v", err)
	}
}
