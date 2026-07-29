// Package errz provides errors and error codes.
package errz

import (
	"errors"
	"sync/atomic"
)

var globalFactory = func() *atomic.Pointer[factory] {
	var p atomic.Pointer[factory]
	p.Store(NewFactory())
	return &p
}()

var (
	// ErrCanceled indicates the operation was canceled.
	ErrCanceled = E(Canceled)
	// ErrUnknown indicates the operation failed for an unknown reason.
	ErrUnknown = E(Unknown)
	// ErrInvalidArgument indicates client supplied an invalid argument.
	ErrInvalidArgument = E(InvalidArgument)
	// ErrDeadlineExceeded indicates deadline expired before operation could complete.
	ErrDeadlineExceeded = E(DeadlineExceeded)
	// ErrNotFound indicates requested entity was not found.
	ErrNotFound = E(NotFound)
	// ErrAlreadyExists indicates client attempted to create an entity that already exists.
	ErrAlreadyExists = E(AlreadyExists)
	// ErrPermissionDenied indicates caller does not have permission.
	ErrPermissionDenied = E(PermissionDenied)
	// ErrResourceExhausted indicates some resource has been exhausted.
	ErrResourceExhausted = E(ResourceExhausted)
	// ErrFailedPrecondition indicates system is not in state required for operation.
	ErrFailedPrecondition = E(FailedPrecondition)
	// ErrAborted indicates operation was aborted by the system.
	ErrAborted = E(Aborted)
	// ErrOutOfRange indicates operation was attempted past the valid range.
	ErrOutOfRange = E(OutOfRange)
	// ErrUnimplemented indicates operation is not implemented.
	ErrUnimplemented = E(Unimplemented)
	// ErrInternal indicates some invariants expected by underlying system have been broken.
	ErrInternal = E(Internal)
	// ErrUnavailable indicates service is currently unavailable.
	ErrUnavailable = E(Unavailable)
	// ErrDataLoss indicates operation resulted in unrecoverable data loss.
	ErrDataLoss = E(DataLoss)
	// ErrUnauthenticated indicates request does not have valid authentication.
	ErrUnauthenticated = E(Unauthenticated)
)

// E creates an error from the given arguments.
func E(args ...any) *Error {
	return globalFactory.Load().E(args...)
}

// F creates an error with a formatted message.
func F(format string, args ...any) *Error {
	return globalFactory.Load().F(format, args...)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Is reports whether any error in err's chain matches any of the targets.
func Is(err error, targets ...error) bool {
	for i := range targets {
		if errors.Is(err, targets[i]) {
			return true
		}
	}
	return false
}

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
func Join(errs ...error) error {
	var errs2 []any
	for _, err := range errs {
		if err != nil {
			errs2 = append(errs2, err)
		}
	}
	if len(errs2) == 0 {
		return nil
	}

	return E(errs2...)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	unwrapper, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return unwrapper.Unwrap() //nolint:wrapcheck // Expected for Unwrap
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
