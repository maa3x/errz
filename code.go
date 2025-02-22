package errz

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Code int

const (
	// Canceled indicates that the operation was Canceled, typically by the caller.
	Canceled Code = 1

	// Unknown indicates that the operation failed for an Unknown reason.
	Unknown Code = 2

	// InvalidArgument indicates that client supplied an invalid argument.
	InvalidArgument Code = 3

	// DeadlineExceeded indicates that deadline expired before the operation could complete.
	DeadlineExceeded Code = 4

	// NotFound indicates that some requested entity (for example, a file or directory) was not found.
	NotFound Code = 5

	// AlreadyExists indicates that client attempted to create an entity (for example, a file or directory) that already exists.
	AlreadyExists Code = 6

	// PermissionDenied indicates that the caller doesn't have permission to execute the specified operation.
	PermissionDenied Code = 7

	// ResourceExhausted indicates that some resource has been exhausted.
	// For example, a per-user quota may be exhausted or the entire file system may be full.
	ResourceExhausted Code = 8

	// FailedPrecondition indicates that the system is not in a state required for the operation's execution.
	FailedPrecondition Code = 9

	// Aborted indicates that operation was Aborted by the system,
	// usually because of a concurrency issue such as a sequencer check failure or transaction abort.
	Aborted Code = 10

	// OutOfRange indicates that the operation was attempted past the valid range (for example, seeking past end-of-file).
	OutOfRange Code = 11

	// Unimplemented indicates that the operation isn't implemented, supported, or enabled in this service.
	Unimplemented Code = 12

	// Internal indicates that some invariants expected by the underlying system have been broken.
	// This code is reserved for serious errors.
	Internal Code = 13

	// Unavailable indicates that the service is currently Unavailable.
	// This is usually temporary, so clients can back off and retry idempotent operations.
	Unavailable Code = 14

	// DataLoss indicates that the operation has resulted in unrecoverable data loss or corruption.
	DataLoss Code = 15

	// Unauthenticated indicates that the request does not have valid authentication credentials for the operation.
	Unauthenticated Code = 16

	minCode = Canceled
	maxCode = Unauthenticated
)

func (c *Code) String() string {
	if c == nil {
		return ""
	}

	switch *c {
	case Canceled:
		return "canceled"
	case Unknown:
		return "unknown"
	case InvalidArgument:
		return "invalid_argument"
	case DeadlineExceeded:
		return "deadline_exceeded"
	case NotFound:
		return "not_found"
	case AlreadyExists:
		return "already_exists"
	case PermissionDenied:
		return "permission_denied"
	case ResourceExhausted:
		return "resource_exhausted"
	case FailedPrecondition:
		return "failed_precondition"
	case Aborted:
		return "aborted"
	case OutOfRange:
		return "out_of_range"
	case Unimplemented:
		return "unimplemented"
	case Internal:
		return "internal"
	case Unavailable:
		return "unavailable"
	case DataLoss:
		return "data_loss"
	case Unauthenticated:
		return "unauthenticated"
	}
	return fmt.Sprintf("code_%d", c)
}

// MarshalText implements [encoding.TextMarshaler].
func (c *Code) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (c *Code) UnmarshalText(data []byte) error {
	dataStr := string(data)
	switch dataStr {
	case "canceled":
		*c = Canceled
		return nil
	case "unknown":
		*c = Unknown
		return nil
	case "invalid_argument":
		*c = InvalidArgument
		return nil
	case "deadline_exceeded":
		*c = DeadlineExceeded
		return nil
	case "not_found":
		*c = NotFound
		return nil
	case "already_exists":
		*c = AlreadyExists
		return nil
	case "permission_denied":
		*c = PermissionDenied
		return nil
	case "resource_exhausted":
		*c = ResourceExhausted
		return nil
	case "failed_precondition":
		*c = FailedPrecondition
		return nil
	case "aborted":
		*c = Aborted
		return nil
	case "out_of_range":
		*c = OutOfRange
		return nil
	case "unimplemented":
		*c = Unimplemented
		return nil
	case "internal":
		*c = Internal
		return nil
	case "unavailable":
		*c = Unavailable
		return nil
	case "data_loss":
		*c = DataLoss
		return nil
	case "unauthenticated":
		*c = Unauthenticated
		return nil
	}
	// Ensure that non-canonical codes round-trip through MarshalText and UnmarshalText.
	if strings.HasPrefix(dataStr, "code_") {
		dataStr = strings.TrimPrefix(dataStr, "code_")
		_code, err := strconv.ParseInt(dataStr, 10 /* base */, 64 /* bitsize */)
		if err == nil && (_code < int64(minCode) || _code > int64(maxCode)) {
			*c = Code(_code)
			return nil
		}
	}
	return fmt.Errorf("invalid code %q", dataStr)
}

// CodeOf returns the error's status code if it is or wraps an [*Error] and [unknown] otherwise.
func CodeOf(err error) Code {
	var e *Error
	ok := errors.As(err, &e)
	if ok {
		return e.code
	}
	return Unknown
}
