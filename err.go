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
	Code        Code
	Message     string
	Location    *location
	Timestamp   *time.Time
	Stackframes stackframes
	Meta        metadata
	errs        []error
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

// Error implements the error interface.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	return e.String()
}

// String returns a string representation of the error.
func (e *Error) String() string {
	if e == nil {
		return ""
	}

	parts := make([]string, 0)
	if e.Timestamp != nil {
		parts = append(parts, e.Timestamp.Format(time.DateTime))
	}
	if e.Code != 0 {
		parts = append(parts, fmt.Sprintf("[%d %s]", e.Code, e.Code.String()))
	}
	if e.Message != "" {
		parts = append(parts, e.Message)
	}
	if e.Location != nil {
		parts = append(parts, e.Location.String())
	}
	if len(e.Meta) > 0 {
		parts = append(parts, e.Meta.String())
	}
	if len(e.Stackframes) > 0 {
		parts = append(parts, e.Stackframes.String())
	}

	var builder strings.Builder
	builder.WriteString(strings.Join(parts, " "))
	e.appendErrs(&builder)

	return builder.String()
}

// IsEmpty reports whether the error contains no information.
func (e *Error) IsEmpty() bool {
	return e == nil || (len(e.errs) == 0 && e.Code == 0 && e.Message == "" && len(e.Meta) == 0)
}

// With adds a metadata key-value pair to the error.
func (e *Error) With(key string, value any) *Error {
	if e == nil {
		return nil
	}

	e.Meta = e.Meta.Add(key, value)
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
	e.Stackframes = make([]runtime.Frame, 0, 8)
	for {
		frame, more := frames.Next()
		e.Stackframes = append(e.Stackframes, frame)
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

	e.Timestamp = new(time.Now())
	return e
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
	e.Location = &location{File: frame.File, Func: frame.Function, Line: frame.Line}

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
