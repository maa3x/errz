package errz

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Error represents a structured error.
type Error struct {
	loc   *location
	ts    *time.Time
	msg   string
	errs  []error
	stack stackframes
	meta  Metadata
	code  Code
}

// With adds a metadata key-value pair to the error.
func (e *Error) With(key string, value any) *Error {
	if e == nil {
		return nil
	}

	e.meta = e.meta.Add(key, value)
	return e
}

// WithMany adds multiple metadata key-value pairs to the error.
func (e *Error) WithMany(pairs ...any) *Error {
	if e == nil {
		return nil
	}

	if len(pairs)%2 != 0 {
		return e.With("errz_internal", "WithMany requires an even number of arguments")
	}

	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return e.With("errz_internal", "WithMany keys must be strings")
		}
		e = e.With(key, pairs[i+1])
	}

	return e
}

// WithLocation adds the caller's location to the error.
func (e *Error) WithLocation() *Error {
	if e == nil {
		return nil
	}
	return e.addLocation(4)
}

// WithTrace adds a stack trace to the error, skipping the specified number of frames.
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

// WithTime adds the current time to the error.
func (e *Error) WithTime() *Error {
	if e == nil {
		return nil
	}

	e.ts = new(time.Now())
	return e
}

// String returns a string representation of the error.
func (e *Error) String() string {
	if e == nil {
		return ""
	}

	parts := make([]string, 0)
	if e.ts != nil {
		parts = append(parts, e.ts.Format(time.DateTime))
	}
	if e.code != 0 {
		parts = append(parts, fmt.Sprintf("[%d %s]", e.code, e.code.String()))
	}
	if e.msg != "" {
		parts = append(parts, e.msg)
	}
	if e.loc != nil {
		parts = append(parts, e.loc.String())
	}
	if len(e.meta) > 0 {
		parts = append(parts, e.meta.String())
	}
	if len(e.stack) > 0 {
		parts = append(parts, e.stack.String())
	}

	var builder strings.Builder
	builder.WriteString(strings.Join(parts, " "))
	e.appendErrs(&builder)

	return builder.String()
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	return e.String()
}

// Unwrap returns the underlying errors.
func (e *Error) Unwrap() []error {
	if e == nil {
		return nil
	}
	return e.errs
}

// Wrap appends an error to the error's chain.
func (e *Error) Wrap(err error) *Error {
	if e == nil {
		return nil
	} else if err == nil {
		return e
	}

	e.errs = append(e.errs, err)
	return e
}

// IsEmpty reports whether the error contains no information.
func (e *Error) IsEmpty() bool {
	return e == nil || (len(e.errs) == 0 && e.code == 0 && e.msg == "" && len(e.meta) == 0)
}

// Code returns the error's code.
func (e *Error) Code() Code {
	if e == nil {
		return 0
	}

	return e.code
}

// Message returns the error's message.
func (e *Error) Message() string {
	if e == nil {
		return ""
	}

	return e.msg
}

// Location returns the caller's location when the error was created.
//
//nolint:revive // unexported-return: location is an internal type
func (e *Error) Location() *location {
	if e == nil {
		return nil
	}

	return e.loc
}

// StackTrace returns the error's stack trace.
func (e *Error) StackTrace() []runtime.Frame {
	if e == nil {
		return nil
	}

	return e.stack
}

// Timestamp returns the time the error was created.
func (e *Error) Timestamp() *time.Time {
	if e == nil {
		return nil
	}

	return e.ts
}

// Meta returns the error's metadata.
func (e *Error) Meta() Metadata {
	if e == nil {
		return nil
	}

	return e.meta
}

func (e *Error) appendErrs(builder *strings.Builder) {
	count := len(e.errs)
	if count == 0 {
		return
	}

	builder.WriteString(": ")
	if count == 1 {
		builder.WriteString(e.errs[0].Error())
		return
	}

	errs := make([]string, len(e.errs))
	for i := range e.errs {
		errs[i] = "(" + e.errs[i].Error() + ")"
	}
	builder.WriteString(strings.Join(errs, " "))
}

func (e *Error) addLocation(skip int) *Error {
	if e == nil {
		return nil
	}

	var pcs [4]uintptr
	frames := runtime.CallersFrames(pcs[:runtime.Callers(skip, pcs[:])])
	frame, _ := frames.Next()
	e.loc = &location{File: frame.File, Func: frame.Function, Line: frame.Line}

	return e
}

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
