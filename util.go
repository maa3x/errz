package errz

import "errors"

var (
	ErrCanceled           = E(Canceled)
	ErrUnknown            = E(Unknown)
	ErrInvalidArgument    = E(InvalidArgument)
	ErrDeadlineExceeded   = E(DeadlineExceeded)
	ErrNotFound           = E(NotFound)
	ErrAlreadyExists      = E(AlreadyExists)
	ErrPermissionDenied   = E(PermissionDenied)
	ErrResourceExhausted  = E(ResourceExhausted)
	ErrFailedPrecondition = E(FailedPrecondition)
	ErrAborted            = E(Aborted)
	ErrOutOfRange         = E(OutOfRange)
	ErrUnimplemented      = E(Unimplemented)
	ErrInternal           = E(Internal)
	ErrUnavailable        = E(Unavailable)
	ErrDataLoss           = E(DataLoss)
	ErrUnauthenticated    = E(Unauthenticated)
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
	for i := range targets {
		if errors.Is(err, targets[i]) {
			return true
		}
	}
	return false
}

func IsNot(err error, targets ...error) bool {
	return !Is(err, targets...)
}
