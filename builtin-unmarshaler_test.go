package bargle

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func parseBuiltinLong[T BuiltinUnmarshalerType](t *testing.T, target *T, args ...string) error {
	p := NewParser()
	p.SetArgs(args...)
	ParseAll(p, Long("value", BuiltinUnmarshaler(target)))
	p.FailIfArgsRemain()
	return p.Err()
}

func TestBuiltinIntegerWidths(t *testing.T) {
	c := qt.New(t)

	var u16 uint16
	c.Assert(parseBuiltinLong(t, &u16, "--value=65535"), qt.IsNil)
	c.Check(u16, qt.Equals, uint16(65535))
	// Out of range for the target's width, rather than silently truncated.
	c.Check(parseBuiltinLong(t, &u16, "--value=65536"), qt.IsNotNil)

	var i16 int16
	c.Assert(parseBuiltinLong(t, &i16, "--value=-32768"), qt.IsNil)
	c.Check(i16, qt.Equals, int16(-32768))
	c.Check(parseBuiltinLong(t, &i16, "--value=32768"), qt.IsNotNil)

	var i32 int32
	c.Assert(parseBuiltinLong(t, &i32, "--value=2147483647"), qt.IsNil)
	c.Check(i32, qt.Equals, int32(2147483647))
	c.Check(parseBuiltinLong(t, &i32, "--value=2147483648"), qt.IsNotNil)

	var i64 int64
	c.Assert(parseBuiltinLong(t, &i64, "--value=-9223372036854775808"), qt.IsNil)
	c.Check(i64, qt.Equals, int64(-9223372036854775808))

	var ui uint
	c.Assert(parseBuiltinLong(t, &ui, "--value=0x10"), qt.IsNil)
	c.Check(ui, qt.Equals, uint(16))
	c.Check(parseBuiltinLong(t, &ui, "--value=-1"), qt.IsNotNil)
}
