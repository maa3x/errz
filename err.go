package errz

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"
)

type Error struct {
	errs  []error
	code  Code
	msg   string
	loc   *location
	ts    *time.Time
	stack []runtime.Frame
	meta  Metadata
}

func (e *Error) With(key string, value any) *Error {
	if e == nil {
		return nil
	}

	e.meta.Add(key, value)
	return e
}

func (e *Error) WithLocation() *Error {
	if e == nil {
		return nil
	}

	var pcs [4]uintptr
	frames := runtime.CallersFrames(pcs[:runtime.Callers(3, pcs[:])])
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
		b.WriteString(fmt.Sprintf("\n   %s\n", e.loc.String()))
	}

	if len(e.meta) > 0 {
		b.WriteString("\n" + e.meta.String())
	}
	if len(e.stack) > 0 {
		buf := new(bytes.Buffer)
		w := tabwriter.NewWriter(buf, 0, 0, 3, ' ', 0)
		for i := range e.stack {
			fmt.Fprintf(w, "\t%s\t%s:%d\n", e.stack[i].Function, e.stack[i].File, e.stack[i].Line)
		}
		w.Flush()
		b.WriteRune('\n')
		b.Write(buf.Bytes())
	}

	if len(e.errs) > 0 {
		b.WriteRune('\n')
		errs := make([]string, len(e.errs))
		for i := range e.errs {
			lines := strings.Split(fmt.Sprint(e.errs[i]), "\n")
			for j := range lines {
				lines[j] = "   " + lines[j]
			}
			errs[i] = strings.Join(lines, "\n")
		}
		b.WriteString(strings.Join(errs, "\n"))
		b.WriteRune('\n')
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
