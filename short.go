package bargle

type ShortParser struct {
	Short    rune
	CanUnary bool
}

func (me ShortParser) Parse(ctx Context) error {
	args := ctx.Args()
	next := []rune(args.Pop())
	if len(next) < 2 {
		return noMatch
	}
	if next[0] != '-' {
		return noMatch
	}
	if next[1] != me.Short {
		return noMatch
	}
	next = next[2:]
	if len(next) == 0 {
		return nil
	}
	if me.CanUnary {
		if next[0] == '=' {
			next = next[1:]
		}
		args.Push(string(next))
	} else {
		args.Push(string(append([]rune{'-'}, next...)))
	}
	return nil
}
