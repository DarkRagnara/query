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
	slice interface{}
}

func NewSliceIter(slice interface{}) SliceIter {
	return SliceIter{slice}
}

func (s SliceIter) Peek() interface{} {
	val := reflect.ValueOf(s.slice)
	return val.Index(0).Interface()
}

func (s SliceIter) IsEmpty() bool {
	val := reflect.ValueOf(s.slice)
	return val.Len() == 0
}

func (s *SliceIter) Next() {
	val := reflect.ValueOf(s.slice)
	val = val.Slice(1, val.Len())
	s.slice = val.Interface()
}
