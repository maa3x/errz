package errz

import "fmt"

type detail struct {
	Key   string
	Value any
}

func (d detail) String() string {
	return d.Key + "\t" + fmt.Sprint(d.Value)
}
