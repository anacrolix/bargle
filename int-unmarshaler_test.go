package bargle

import (
	"strconv"
	"testing"

	g "github.com/anacrolix/generics"
	qt "github.com/frankban/quicktest"
)

type intMarshalCase struct {
	arg string
	i   g.Option[int64]
}

func goodCase(arg string, value int64) intMarshalCase {
	return intMarshalCase{arg, g.Some[int64](value)}
}

func badCase(arg string) intMarshalCase {
	return intMarshalCase{arg, g.None[int64]()}
}

var intMarshalCases = []struct {
	arg string
	i   g.Option[int64]
}{
	// Check the float parsing with the nudge toward zero.
	goodCase("0e0", 0),

	// Typical values with scientific literal and not.
	goodCase("10e3", 10e3),
	goodCase("10e3", 10000),

	// Decimal point but still valid integer
	goodCase("-4.123e9", -4_123_000_000),
	// float64 has around 16 decimal digits of precision. This value can't be represented by a float64.
	badCase("-12345678901234.567e3"),
	// This one can.
	goodCase("-12345678901234.56e2", -1234567890123456),
}

func TestUnmarshalInt(t *testing.T) {
	c := qt.New(t)
	for _, _case := range intMarshalCases {
		ui, err := unmarshalInt[int64, int64](_case.arg, strconv.ParseInt, 64)
		if _case.i.Ok {
			c.Assert(err, qt.IsNil)
			c.Check(ui, qt.Equals, _case.i.Value)
		} else {
			c.Assert(err, qt.IsNotNil)
		}
	}
}

func TestFloatFormatting(t *testing.T) {
	c := qt.New(t)
	c.Check(strconv.FormatFloat(10000, 'e', -1, 64), qt.Equals,
		"1e+04")
	// This one exceeds precision, and still doesn't trigger exponent or decimal points despite 'f'.
	c.Check(strconv.FormatFloat(-12345678901234567, 'f', -1, 64), qt.Equals,
		"-12345678901234568")
}

func FuzzIntUnmarshalling(f *testing.F) {
	for _, _case := range intMarshalCases {
		if _case.i.Ok {
			f.Add(_case.i.Value)
		}
	}
	f.Fuzz(func(t *testing.T, i int64) {
		fs := strconv.FormatFloat(float64(i), 'e', -1, 64)
		//if !strings.Contains(fs, "e") {
		//	t.SkipNow()
		//}
		ui, err := unmarshalInt[int64, int64](fs, strconv.ParseInt, 64)
		if err != nil {
			t.Fatal(err)
		}
		if ui != i {
			t.FailNow()
		}
	})
}
