package errz

import (
	"runtime"
	"strings"
)

type stackframes []runtime.Frame

func (s stackframes) String() string {
	if len(s) == 0 {
		return ""
	}

	var frames []string
	for i := range s {
		loc := &location{File: s[i].File, Func: s[i].Function, Line: s[i].Line}
		frames = append(frames, loc.String())
	}
	return strings.Join(frames, " ")
}
