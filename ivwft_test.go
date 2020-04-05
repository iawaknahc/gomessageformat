package messageformat

import (
	"testing"
)

func TestIVWFT(t_ *testing.T) {
	test := func(input interface{}, ii, vv, ww, ff, tt int) {
		i, v, w, f, t, err := IVWFT(input)
		if err != nil {
			t_.Errorf("err: %v\n", err)
		} else {
			if i != ii {
				t_.Errorf("%d != %d\n", i, ii)
			}
			if v != vv {
				t_.Errorf("%d != %d\n", v, vv)
			}
			if w != ww {
				t_.Errorf("%d != %d\n", w, ww)
			}
			if f != ff {
				t_.Errorf("%d != %d\n", f, ff)
			}
			if t != tt {
				t_.Errorf("%d != %d\n", t, tt)
			}
		}
	}

	test(1, 1, 0, 0, 0, 0)
	test(1.0, 1, 0, 0, 0, 0)
	test("1", 1, 0, 0, 0, 0)
	test("01", 1, 0, 0, 0, 0)

	test("1.0", 1, 1, 0, 0, 0)
	test("1.00", 1, 2, 0, 0, 0)

	test(1.3, 1, 1, 1, 3, 3)
	test("1.3", 1, 1, 1, 3, 3)

	test("1.30", 1, 2, 1, 30, 3)

	test(1.03, 1, 2, 2, 3, 3)
	test("1.03", 1, 2, 2, 3, 3)

	test(1.23, 1, 2, 2, 23, 23)
	test("1.230", 1, 3, 2, 230, 23)

	test(1234, 1234, 0, 0, 0, 0)
	test(1234.0, 1234, 0, 0, 0, 0)
}
