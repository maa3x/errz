package errz

import (
	"strconv"
)

type location struct {
	File string
	Func string
	Line int
}

func (l *location) String() string {
	if l == nil {
		return ""
	}

	return "[" + l.Func + " " + l.File + ":" + strconv.Itoa(l.Line) + "] "
}
