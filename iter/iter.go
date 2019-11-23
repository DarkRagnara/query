package iter

import (
	"reflect"
)

type ConstIter interface {
	Peek() interface{}
	IsEmpty() bool
}

type Iter interface {
	ConstIter
	Next()
}

type SliceIter struct {
	Slice interface{}
}

func (s SliceIter) Peek() interface{} {
	val := reflect.ValueOf(s.Slice)
	return val.Index(0).Interface()
}

func (s SliceIter) IsEmpty() bool {
	val := reflect.ValueOf(s.Slice)
	return val.Len() == 0
}

func (s *SliceIter) Next() {
	val := reflect.ValueOf(s.Slice)
	val = val.Slice(1, val.Len())
	s.Slice = val.Interface()
}
