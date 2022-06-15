package bargle

type args struct {
	ss []string
}

func NewArgs(ss []string) Args {
	args := &args{}
	for i := range ss {
		args.Push(ss[len(ss)-1-i])
	}
	return args
}

func (me *args) Push(s string) {
	me.ss = append(me.ss, s)
}

func (me *args) Pop() string {
	i := len(me.ss) - 1
	s := me.ss[i]
	me.ss = me.ss[:i]
	return s
}

func (me *args) Len() int {
	return len(me.ss)
}

func (me *args) Clone() Args {
	return &args{
		ss: append(make([]string, 0, len(me.ss)), me.ss...),
	}
}

type Args interface {
	Push(string)
	Pop() string
	Len() int
	Clone() Args
}
