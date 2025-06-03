package errz

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const (
	indent = "  "
	sep    = "\n"
)

type Error struct {
	errs  []error
	code  Code
	msg   string
	loc   *location
	ts    *time.Time
	stack stackframes
	meta  Metadata
}

func (e *Error) With(key string, value any) *Error {
	if e == nil {
		return nil
	}

	e.meta = e.meta.Add(key, value)
	return e
}

func (e *Error) WithLocation() *Error {
	if e == nil {
		return nil
	}
	return e.withLocation(4)
}

func (e *Error) withLocation(skip int) *Error {
	if e == nil {
		return nil
	}

	var pcs [4]uintptr
	frames := runtime.CallersFrames(pcs[:runtime.Callers(skip, pcs[:])])
	frame, _ := frames.Next()
	e.loc = &location{File: frame.File, Func: frame.Function, Line: frame.Line}
	return e
}

func (e *Error) WithTrace(skip int) *Error {
	if e == nil {
		return nil
	}

	pcs := make([]uintptr, 32)
	frames := runtime.CallersFrames(pcs[:runtime.Callers(3+skip, pcs)])
	e.stack = make([]runtime.Frame, 0, 8)
	for {
		frame, more := frames.Next()
		e.stack = append(e.stack, frame)
		if !more {
			break
		}
	}

	return e
}

func (e *Error) WithTime() *Error {
	if e == nil {
		return nil
	}

	now := time.Now()
	e.ts = &now
	return e
}

func (e *Error) String() string {
	if e == nil {
		return ""
	}

	var b strings.Builder
	if e.ts != nil {
		b.WriteString(e.ts.Format(time.DateTime) + " ")
	}
	if e.code != 0 {
		b.WriteString(fmt.Sprintf("[%d %s]", e.code, e.code.String()))
	}
	if e.msg != "" {
		b.WriteString(e.msg + ": ")
	}
	if e.loc != nil {
		b.WriteString(e.loc.String())
	}
	if len(e.meta) > 0 {
		b.WriteString(e.meta.String())
	}
	if len(e.stack) > 0 {
		b.WriteString(e.stack.String())
	}

	if len(e.errs) > 0 {
		errs := make([]string, len(e.errs))
		for i := range e.errs {
			errs[i] = "(" + e.errs[i].Error() + ")"
		}
		b.WriteString(strings.Join(errs, " "))
	}
	return b.String()
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	return e.String()
}

func (e *Error) Unwrap() []error {
	if e == nil {
		return nil
	}

	return e.errs
}

func (e *Error) Wrap(err error) *Error {
	if e == nil {
		return nil
	}

	if err == nil {
		return e
	}

	e.errs = append(e.errs, err)
	return e
}

func (e *Error) Is(targets ...error) bool {
	if e == nil {
		return false
	}

	return Is(e, targets...)
}

func (e *Error) As(err error) bool {
	if e == nil {
		return false
	}

	return As(e, err)
}

func (e *Error) IsEmpty() bool {
	return e == nil || (len(e.errs) == 0 && e.code == 0 && e.msg == "" && len(e.meta) == 0)
}

func (e *Error) Code() Code {
	if e == nil {
		return 0
	}

	return e.code
}

func (e *Error) Message() string {
	if e == nil {
		return ""
	}

	return e.msg
}

func (e *Error) Location() *location {
	if e == nil {
		return nil
	}

	return e.loc
}

func (e *Error) StackTrace() []runtime.Frame {
	if e == nil {
		return nil
	}

	return e.stack
}

func (e *Error) Timestamp() *time.Time {
	if e == nil {
		return nil
	}

	return e.ts
}

func (e *Error) Meta() Metadata {
	if e == nil {
		return nil
	}

	return e.meta
}
