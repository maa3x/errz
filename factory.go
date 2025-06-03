package errz

import (
	"fmt"
	"strings"
)

var DefaultFactory = Factory()

func F(format string, args ...any) *Error {
	return DefaultFactory.F(format, args...)
}

func E(in ...any) *Error {
	return DefaultFactory.E(in...)
}

func If(err error, in ...any) *Error {
	if err == nil {
		return nil
	}

	return E(append([]any{err}, in...)...)
}

type factory struct {
	stacktrace bool
	location   bool
	timestamp  bool
	OnError    func(*Error)
}

func Factory() *factory {
	return &factory{location: true}
}

func (f *factory) StackTrace(v bool) *factory {
	f.stacktrace = v
	return f
}

func (f *factory) Location(v bool) *factory {
	f.location = v
	return f
}

func (f *factory) Timestamp(v bool) *factory {
	f.timestamp = v
	return f
}

func (f *factory) E(in ...any) *Error {
	err := &Error{}
	for i := range in {
		if in[i] == nil {
			continue
		}

		switch v := in[i].(type) {
		case Code:
			err.code = v
		case Error:
			if !v.IsEmpty() {
				err.errs = append(err.errs, &v)
			}
		case *Error:
			if v != nil && !v.IsEmpty() {
				err.errs = append(err.errs, v)
			}
		case error:
			if v != nil {
				err.errs = append(err.errs, v)
			}
		case string:
			err.msg = strings.TrimSpace(v)
		case fmt.Stringer:
			if v != nil {
				err.msg = strings.TrimSpace(v.String())
			}
		default:
			if v != nil {
				err.meta = append(err.meta, detail{Value: v})
			}
		}
	}

	if err.IsEmpty() {
		return nil
	}

	if f.stacktrace {
		err.WithTrace(1)
	}
	if f.location {
		err.withLocation(4)
	}
	if f.timestamp {
		err.WithTime()
	}
	return err
}

func (f *factory) F(format string, args ...any) *Error {
	return E(fmt.Errorf(format, args...))
}
