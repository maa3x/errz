package errz_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/maa3x/errz"
)

func TestIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		targets  []error
		expected bool
	}{
		{
			name:     "nil error matches nothing",
			err:      nil,
			targets:  []error{errz.ErrNotFound, errz.ErrUnknown},
			expected: false,
		},
		{
			name:     "error matches itself",
			err:      errz.ErrNotFound,
			targets:  []error{errz.ErrNotFound},
			expected: true,
		},
		{
			name:     "error matches one of multiple targets",
			err:      errz.ErrNotFound,
			targets:  []error{errz.ErrUnknown, errz.ErrNotFound, errz.ErrInternal},
			expected: true,
		},
		{
			name:     "error doesn't match different errors",
			err:      errz.ErrNotFound,
			targets:  []error{errz.ErrUnknown, errz.ErrInternal},
			expected: false,
		},
		{
			name:     "wrapped error matches its target",
			err:      errz.Join(errz.ErrNotFound),
			targets:  []error{errz.ErrNotFound},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errz.Is(tt.err, tt.targets...); got != tt.expected {
				t.Errorf("Is() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNested(t *testing.T) {
	e1 := errors.New("error 1")
	e2 := errz.E("error 2", e1).WithLocation()
	if !errz.Is(e2, e1) {
		t.Error("error 1 should be nested in error 2")
	}

	e3 := errz.E("error 3", e1, e2).With("meta", "data").With("second", 2)
	if !errz.Is(e3, e1) {
		t.Error("error 1 should be nested in error 3")
	}

	e4 := fmt.Errorf("error 4: %w", e3)
	if !errz.Is(e4, e1) {
		t.Error("error 1 should be nested in error 4")
	}
	if !errz.Is(e4, e2) {
		t.Error("error 1 should be nested in error 2")
	}
	if !errors.Is(e4, e3) {
		t.Error("error 1 should be nested in error 3")
	}
}
