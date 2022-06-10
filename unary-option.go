package bargle

type unaryOption[T any] struct {
	Value       T
	Unmarshaler UnaryUnmarshaler[T]
	Longs       []string
	Shorts      []rune
}

func UnaryOption[T any](target UnaryUnmarshaler[T]) *unaryOption[T] {
	return &unaryOption[T]{Unmarshaler: target}
}

func (me *unaryOption[T]) switchForms() (ret []string) {
	for _, l := range me.Longs {
		ret = append(ret, "--"+l)
	}
	for _, s := range me.Shorts {
		ret = append(ret, "-"+string(s))
	}
	return
}

func (me *unaryOption[T]) Help(f HelpFormatter) {
	ph := ParamHelp{
		Forms: me.switchForms(),
	}
	me.Unmarshaler.Help(&ph)
	f.AddOption(ph)
}

func (me *unaryOption[T]) AddLong(long string) *unaryOption[T] {
	me.Longs = append(me.Longs, long)
	return me
}

func (me *unaryOption[T]) AddShort(short rune) *unaryOption[T] {
	me.Shorts = append(me.Shorts, short)
	return me
}

func (me *unaryOption[T]) Parse(ctx Context) error {
	if !me.matchSwitch(ctx) {
		return noMatch
	}
	return me.Unmarshaler.Unmarshal(ctx.Args().Pop(), &me.Value)
}

func (me unaryOption[T]) matchSwitch(ctx Context) bool {
	for _, l := range me.Longs {
		if ctx.Match(LongParser{Long: l, CanUnary: true}) {
			return true
		}
	}
	for _, s := range me.Shorts {
		if ctx.Match(ShortParser{Short: s, CanUnary: true}) {
			return true
		}
	}
	return false
}
