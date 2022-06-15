package bargle

type ShortParser struct {
	Short    rune
	CanUnary bool
	Prefix   rune
	gotValue bool
}

func (me ShortParser) GotValue() bool {
	return me.gotValue
}

func (me *ShortParser) Parse(ctx Context) error {
	args := ctx.Args()
	if args.Len() == 0 {
		return noMatch
	}
	next := []rune(args.Pop())
	if len(next) < 2 {
		return noMatch
	}
	if me.Prefix == 0 {
		me.Prefix = '-'
	}
	if next[0] != me.Prefix {
		return noMatch
	}
	if next[1] != me.Short {
		return noMatch
	}
	me.gotValue = false
	next = next[2:]
	if len(next) == 0 {
		return nil
	}
	if me.CanUnary {
		if next[0] == '=' {
			next = next[1:]
		}
		me.gotValue = true
		args.Push(string(next))
	} else {
		args.Push(string(append([]rune{me.Prefix}, next...)))
	}
	return nil
}

func (me ShortParser) Help(f *ParamHelp) {
	f.Forms = append(f.Forms, "-"+string(me.Short))
}
