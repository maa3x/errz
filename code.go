package errz

import (
	"fmt"
	"strconv"
	"strings"
)

// Code represents an error code.
type Code uint32

const (
	// OK is returned on success.
	OK Code = 0

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

	// AlreadyExists indicates that client attempted to create an entity (e.g, a file or directory) that already exists.
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

	minCode = OK
	maxCode = Unauthenticated
)

var codeStrings = [...]string{
	OK:                 "ok",
	Canceled:           "canceled",
	Unknown:            "unknown",
	InvalidArgument:    "invalid_argument",
	DeadlineExceeded:   "deadline_exceeded",
	NotFound:           "not_found",
	AlreadyExists:      "already_exists",
	PermissionDenied:   "permission_denied",
	ResourceExhausted:  "resource_exhausted",
	FailedPrecondition: "failed_precondition",
	Aborted:            "aborted",
	OutOfRange:         "out_of_range",
	Unimplemented:      "unimplemented",
	Internal:           "internal",
	Unavailable:        "unavailable",
	DataLoss:           "data_loss",
	Unauthenticated:    "unauthenticated",
}

var codeValues = map[string]Code{
	"ok":                  OK,
	"canceled":            Canceled,
	"unknown":             Unknown,
	"invalid_argument":    InvalidArgument,
	"deadline_exceeded":   DeadlineExceeded,
	"not_found":           NotFound,
	"already_exists":      AlreadyExists,
	"permission_denied":   PermissionDenied,
	"resource_exhausted":  ResourceExhausted,
	"failed_precondition": FailedPrecondition,
	"aborted":             Aborted,
	"out_of_range":        OutOfRange,
	"unimplemented":       Unimplemented,
	"internal":            Internal,
	"unavailable":         Unavailable,
	"data_loss":           DataLoss,
	"unauthenticated":     Unauthenticated,
}

func (c *Code) String() string {
	if c == nil {
		return ""
	}

	if *c <= maxCode {
		return codeStrings[*c]
	}
	return fmt.Sprintf("code_%d", *c)
}

// MarshalText implements [encoding.TextMarshaler].
func (c *Code) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (c *Code) UnmarshalText(data []byte) error {
	dataStr := string(data)
	if code, ok := codeValues[dataStr]; ok {
		*c = code
		return nil
	}

	// Ensure that non-canonical codes round-trip through MarshalText and UnmarshalText.
	if after, ok := strings.CutPrefix(dataStr, "code_"); ok {
		code, err := strconv.ParseUint(after, 10 /* base */, 32 /* bitsize */)
		if err == nil && code > uint64(maxCode) {
			*c = Code(code)
			return nil
		}
	}

	return ErrInvalidArgument
}

// CodeOf returns the error's status code if it is or wraps an [*Error] and [Unknown] otherwise.
func CodeOf(err error) Code {
	var e *Error
	if ok := As(err, &e); ok {
		return e.code
	}
	return Unknown
}
