package errz

import (
	"fmt"
	"strings"
)

type detail struct {
	Value any
	Key   string
}

func (d detail) String() string {
	return d.Key + ": " + fmt.Sprint(d.Value)
}

// metadata is a collection of metadata details.
type metadata []detail

// Add adds a key-value pair to the metadata.
func (m metadata) Add(key string, value any) metadata {
	m = append(m, detail{Key: key, Value: value})
	return m
}

// String returns a string representation of the metadata.
func (m metadata) String() string {
	if len(m) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteByte('{')
	for i := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(m[i].String())
	}
	builder.WriteString("} ")

	return builder.String()
}

// Get returns all values associated with the given key.
func (m metadata) Get(key string) []any {
	if m == nil {
		return nil
	}

	var values []any
	for _, d := range m {
		if d.Key == key {
			values = append(values, d.Value)
		}
	}

	return values
}

// Has reports whether the metadata contains the given key.
func (m metadata) Has(key string) bool {
	if m == nil {
		return false
	}

	for _, d := range m {
		if d.Key == key {
			return true
		}
	}

	return false
}
