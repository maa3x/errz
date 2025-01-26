package errz

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type code int

const (
	// canceled indicates that the operation was canceled, typically by the caller.
	canceled code = 1

	// unknown indicates that the operation failed for an unknown reason.
	unknown code = 2

	// invalidArgument indicates that client supplied an invalid argument.
	invalidArgument code = 3

	// deadlineExceeded indicates that deadline expired before the operation could complete.
	deadlineExceeded code = 4

	// notFound indicates that some requested entity (for example, a file or directory) was not found.
	notFound code = 5

	// alreadyExists indicates that client attempted to create an entity (for example, a file or directory) that already exists.
	alreadyExists code = 6

	// permissionDenied indicates that the caller doesn't have permission to execute the specified operation.
	permissionDenied code = 7

	// resourceExhausted indicates that some resource has been exhausted.
	// For example, a per-user quota may be exhausted or the entire file system may be full.
	resourceExhausted code = 8

	// failedPrecondition indicates that the system is not in a state required for the operation's execution.
	failedPrecondition code = 9

	// aborted indicates that operation was aborted by the system,
	// usually because of a concurrency issue such as a sequencer check failure or transaction abort.
	aborted code = 10

	// outOfRange indicates that the operation was attempted past the valid range (for example, seeking past end-of-file).
	outOfRange code = 11

	// unimplemented indicates that the operation isn't implemented, supported, or enabled in this service.
	unimplemented code = 12

	// internal indicates that some invariants expected by the underlying system have been broken.
	// This code is reserved for serious errors.
	internal code = 13

	// unavailable indicates that the service is currently unavailable.
	// This is usually temporary, so clients can back off and retry idempotent operations.
	unavailable code = 14

	// dataLoss indicates that the operation has resulted in unrecoverable data loss or corruption.
	dataLoss code = 15

	// unauthenticated indicates that the request does not have valid authentication credentials for the operation.
	unauthenticated code = 16

	minCode = canceled
	maxCode = unauthenticated
)

func (c *code) String() string {
	if c == nil {
		return ""
	}

	switch *c {
	case canceled:
		return "canceled"
	case unknown:
		return "unknown"
	case invalidArgument:
		return "invalid_argument"
	case deadlineExceeded:
		return "deadline_exceeded"
	case notFound:
		return "not_found"
	case alreadyExists:
		return "already_exists"
	case permissionDenied:
		return "permission_denied"
	case resourceExhausted:
		return "resource_exhausted"
	case failedPrecondition:
		return "failed_precondition"
	case aborted:
		return "aborted"
	case outOfRange:
		return "out_of_range"
	case unimplemented:
		return "unimplemented"
	case internal:
		return "internal"
	case unavailable:
		return "unavailable"
	case dataLoss:
		return "data_loss"
	case unauthenticated:
		return "unauthenticated"
	}
	return fmt.Sprintf("code_%d", c)
}

// MarshalText implements [encoding.TextMarshaler].
func (c *code) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (c *code) UnmarshalText(data []byte) error {
	dataStr := string(data)
	switch dataStr {
	case "canceled":
		*c = canceled
		return nil
	case "unknown":
		*c = unknown
		return nil
	case "invalid_argument":
		*c = invalidArgument
		return nil
	case "deadline_exceeded":
		*c = deadlineExceeded
		return nil
	case "not_found":
		*c = notFound
		return nil
	case "already_exists":
		*c = alreadyExists
		return nil
	case "permission_denied":
		*c = permissionDenied
		return nil
	case "resource_exhausted":
		*c = resourceExhausted
		return nil
	case "failed_precondition":
		*c = failedPrecondition
		return nil
	case "aborted":
		*c = aborted
		return nil
	case "out_of_range":
		*c = outOfRange
		return nil
	case "unimplemented":
		*c = unimplemented
		return nil
	case "internal":
		*c = internal
		return nil
	case "unavailable":
		*c = unavailable
		return nil
	case "data_loss":
		*c = dataLoss
		return nil
	case "unauthenticated":
		*c = unauthenticated
		return nil
	}
	// Ensure that non-canonical codes round-trip through MarshalText and UnmarshalText.
	if strings.HasPrefix(dataStr, "code_") {
		dataStr = strings.TrimPrefix(dataStr, "code_")
		_code, err := strconv.ParseInt(dataStr, 10 /* base */, 64 /* bitsize */)
		if err == nil && (_code < int64(minCode) || _code > int64(maxCode)) {
			*c = code(_code)
			return nil
		}
	}
	return fmt.Errorf("invalid code %q", dataStr)
}

// CodeOf returns the error's status code if it is or wraps an [*Error] and [unknown] otherwise.
func CodeOf(err error) code {
	var e *Error
	ok := errors.As(err, &e)
	if ok {
		return e.code
	}
	return unknown
}
