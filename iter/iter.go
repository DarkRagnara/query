//Package iter implements an iterator interface as used by query.
//
//An iterator as implemented in this package is a type that may point to
//an element that can be looked at via Peek. Peek may only be called
//if the iterator is not empty, i.e. after the last element it iterates
//over. This can be checked via IsEmpty. Move the iterator to its next
//element (or the end) via calls to Next.
//
//Next and Peek are allowed to panic if IsEmpty returns true.
package iter

import (
	"reflect"
)

//ConstIter contains all methods for an iterator that is not moved
//to the next element. As such, it is not very useful except as a
//parameter to a function that is not supposed to change the iterator.
type ConstIter interface {
	Peek() interface{}
	IsEmpty() bool
}

//Iter is an iterator that can be moved to the next element. Otherwise
//it is the same as a ConstIter.
type Iter interface {
	ConstIter
	Next()
}

//SliceIter is a helper type that takes any slice type and returns an
//iterator pointing to the first element of the slice.
type SliceIter struct {
	slice interface{}
}

//NewSliceIter takes any slice type parameter and returns a SliceIter.
//
//Calling NewSliceIter with other values is not allowed and might cause
//a panic.
func NewSliceIter(slice interface{}) SliceIter {
	return SliceIter{slice}
}

//Peek returns the current slice element.
func (s SliceIter) Peek() interface{} {
	val := reflect.ValueOf(s.slice)
	return val.Index(0).Interface()
}

//IsEmpty returns whether the iterator points to a valid slice element.
func (s SliceIter) IsEmpty() bool {
	val := reflect.ValueOf(s.slice)
	return val.Len() == 0
}

//Next moves the iterator to the next element of the slice. If the iterator
//points to the last slice element, the iterator is empty after this call.
func (s *SliceIter) Next() {
	val := reflect.ValueOf(s.slice)
	val = val.Slice(1, val.Len())
	s.slice = val.Interface()
}
