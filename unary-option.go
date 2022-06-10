package bargle

type UnaryOption[T any] struct {
	Value       T
	Unmarshaler UnaryUnmarshaler[T]
	Longs       []string
	Shorts      []rune
}

func NewUnaryOption[T builtinUnaryUnmarshalTarget](target *T) *UnaryOption[T] {
	return &UnaryOption[T]{Unmarshaler: BuiltinUnaryUnmarshaler[T]{}}
}

func (me *UnaryOption[T]) switchForms() (ret []string) {
	for _, l := range me.Longs {
		ret = append(ret, "--"+l)
	}
	for _, s := range me.Shorts {
		ret = append(ret, "-"+string(s))
	}
	return
}

func (me *UnaryOption[T]) initUnmarshaler() {
	if me.Unmarshaler != nil {
		return
	}
	me.Unmarshaler = BuiltinUnaryUnmarshaler[T]{}
}

func (me *UnaryOption[T]) Help(f HelpFormatter) {
	ph := ParamHelp{
		Forms: me.switchForms(),
	}
	me.initUnmarshaler()
	me.Unmarshaler.Help(&ph)
	f.AddOption(ph)
}

func (me *UnaryOption[T]) AddLong(long string) *UnaryOption[T] {
	me.Longs = append(me.Longs, long)
	return me
}

func (me *UnaryOption[T]) AddShort(short rune) *UnaryOption[T] {
	me.Shorts = append(me.Shorts, short)
	return me
}

func (me *UnaryOption[T]) Parse(ctx Context) error {
	if !me.matchSwitch(ctx) {
		return noMatch
	}
	me.initUnmarshaler()
	return me.Unmarshaler.UnaryUnmarshal(ctx.Args().Pop(), &me.Value)
}

func (me UnaryOption[T]) matchSwitch(ctx Context) bool {
	for _, l := range me.Longs {
		if ctx.Match(&LongParser{Long: l, CanUnary: true}) {
			return true
		}
	}
	for _, s := range me.Shorts {
		if ctx.Match(&ShortParser{Short: s, CanUnary: true}) {
			return true
		}
	}
	return false
}
