package errz

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

type detail struct {
	Key   string
	Value any
}

func (d detail) String() string {
	return d.Key + "\t" + fmt.Sprint(d.Value)
}

type Metadata []detail

func (m Metadata) Add(key string, value any) Metadata {
	if m == nil {
		return nil
	}

	m = append(m, detail{Key: key, Value: value})
	return m
}

func (m Metadata) String() string {
	if m == nil {
		return ""
	}

	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 3, ' ', 0)
	for i := range m {
		fmt.Fprintf(w, "\t%s\n", m[i].String())
	}
	w.Flush()
	return buf.String()
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
