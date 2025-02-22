package errz

import "fmt"

type location struct {
	File string
	Func string
	Line int
}

func (l *location) String() string {
	if l == nil {
		return ""
	}

	return fmt.Sprintf("%s   %s:%d\n", l.Func, l.File, l.Line)
}
