package errz

import "errors"

var (
	ErrCanceled           = E(canceled)
	ErrUnknown            = E(unknown)
	ErrInvalidArgument    = E(invalidArgument)
	ErrDeadlineExceeded   = E(deadlineExceeded)
	ErrNotFound           = E(notFound)
	ErrAlreadyExists      = E(alreadyExists)
	ErrPermissionDenied   = E(permissionDenied)
	ErrResourceExhausted  = E(resourceExhausted)
	ErrFailedPrecondition = E(failedPrecondition)
	ErrAborted            = E(aborted)
	ErrOutOfRange         = E(outOfRange)
	ErrUnimplemented      = E(unimplemented)
	ErrInternal           = E(internal)
	ErrUnavailable        = E(unavailable)
	ErrDataLoss           = E(dataLoss)
	ErrUnauthenticated    = E(unauthenticated)
)

func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

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

func As(err error, target any) bool {
	return errors.As(err, target)
}

func NotAs(err error, target any) bool {
	return !As(err, target)
}

func Is(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func IsNot(err error, targets ...error) bool {
	return !Is(err, targets...)
}
