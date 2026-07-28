// Package errz provides errors and error codes.
package errz

import "sync/atomic"

var globalFactory atomic.Pointer[factory]

func init() {
	f := NewFactory()
	globalFactory.Store(f)
}

// E creates an error from the given arguments.
func E(args ...any) *Error {
	return globalFactory.Load().E(args...)
}

// F creates an error with a formatted message.
func F(format string, args ...any) *Error {
	return globalFactory.Load().F(format, args...)
}

// If creates an error from the given arguments if err is not nil.
func If(err error, args ...any) *Error {
	if err == nil {
		return nil
	}

	return E(append([]any{err}, args...)...)
}

// Of returns the error as an *Error and whether the conversion succeeded.
func Of(err error) (*Error, bool) {
	var e *Error
	ok := As(err, &e)
	return e, ok
}

// ReplaceGlobal changes the global error factory.
func ReplaceGlobal(f *factory) error {
	if f == nil {
		return ErrInvalidArgument
	}

	globalFactory.Store(f)
	return nil
}
