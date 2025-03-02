package errz

import (
	"errors"
	"fmt"
	"testing"
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
			targets:  []error{ErrNotFound, ErrUnknown},
			expected: false,
		},
		{
			name:     "error matches itself",
			err:      ErrNotFound,
			targets:  []error{ErrNotFound},
			expected: true,
		},
		{
			name:     "error matches one of multiple targets",
			err:      ErrNotFound,
			targets:  []error{ErrUnknown, ErrNotFound, ErrInternal},
			expected: true,
		},
		{
			name:     "error doesn't match different errors",
			err:      ErrNotFound,
			targets:  []error{ErrUnknown, ErrInternal},
			expected: false,
		},
		{
			name:     "wrapped error matches its target",
			err:      Join(ErrNotFound),
			targets:  []error{ErrNotFound},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.err, tt.targets...); got != tt.expected {
				t.Errorf("Is() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNested(t *testing.T) {
	e1 := errors.New("error 1")
	e2 := E("error 2", e1).WithLocation()
	if !Is(e2, e1) {
		t.Errorf("error 1 should be nested in error 2")
	}

	e3 := E("error 3", e1, e2).With("meta", "data").With("second", 2)
	if !Is(e3, e1) {
		t.Errorf("error 1 should be nested in error 3")
	}

	e4 := fmt.Errorf("error 4: %w", e3)
	e5 := E("error 5", e1, e2, e3, e4).WithTrace(0)
	fmt.Println(e5)
	if !Is(e4, e1) {
		t.Errorf("error 1 should be nested in error 4")
	}
	if !Is(e4, e2) {
		t.Errorf("error 1 should be nested in error 2")
	}
	if !errors.Is(e4, e3) {
		t.Errorf("error 1 should be nested in error 3")
	}
}
