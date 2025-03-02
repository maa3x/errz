package errz

import (
	"fmt"
)

type detail struct {
	Key   string
	Value any
}

func (d detail) String() string {
	return d.Key + ": " + fmt.Sprint(d.Value)
}

type Metadata []detail

func (m Metadata) Add(key string, value any) Metadata {
	m = append(m, detail{Key: key, Value: value})
	return m
}

func (m Metadata) String() string {
	if len(m) == 0 {
		return ""
	}

	s := "{"
	for i := range m {
		if i > 0 {
			s += ", "
		}
		s += m[i].String()
	}
	return s + "} "
}

func (m Metadata) Get(key string) []any {
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

func (m Metadata) Has(key string) bool {
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
