package errz

import (
	"fmt"
	"strconv"
	"strings"
)

type factory struct {
	OnError    func(*Error)
	stacktrace bool
	location   bool
	timestamp  bool
}

// NewFactory creates a new error factory.
func NewFactory() *factory {
	return &factory{}
}

// E creates an error from the given arguments.
func (f *factory) E(args ...any) *Error {
	err := &Error{}
	f.applyArgs(err, args)
	if err.IsEmpty() {
		return nil
	}

	if f.stacktrace {
		err.WithTrace(2)
	}
	if f.location {
		err.addLocation(5)
	}
	if f.timestamp {
		err.WithTime()
	}

	return err
}

// F creates an error with a formatted message.
func (f *factory) F(format string, args ...any) *Error {
	return f.E(fmt.Errorf(format, args...)) //nolint:err113 // intended
}

// StackTrace enables or disables stack traces.
func (f *factory) StackTrace(enable bool) *factory {
	f.stacktrace = enable
	return f
}

// Location enables or disables caller location.
func (f *factory) Location(enable bool) *factory {
	f.location = enable
	return f
}

// Timestamp enables or disables timestamps.
func (f *factory) Timestamp(enable bool) *factory {
	f.timestamp = enable
	return f
}

//nolint:revive,cyclop // cognitive-complexity: switch logic is flat
func (f *factory) applyArgs(err *Error, args []any) {
	for i := range args {
		if args[i] == nil {
			continue
		}

		switch typedVal := args[i].(type) {
		case Code:
			err.Code = typedVal
		case Error:
			if !typedVal.IsEmpty() {
				err.errs = append(err.errs, &typedVal)
			}
		case *Error:
			if typedVal != nil && !typedVal.IsEmpty() {
				err.errs = append(err.errs, typedVal)
			}
		case error:
			if typedVal != nil {
				err.errs = append(err.errs, typedVal)
			}
		case string:
			err.Message = strings.TrimSpace(typedVal)
		case fmt.Stringer:
			if typedVal != nil {
				err.Message = strings.TrimSpace(typedVal.String())
			}
		default:
			if typedVal != nil {
				err.Meta = append(err.Meta, detail{Key: "arg" + strconv.Itoa(i), Value: typedVal})
			}
		}
	}
}
