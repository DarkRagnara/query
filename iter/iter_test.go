package iter_test

import (
	. "bitbucket.org/ragnara/query/iter"
	"testing"
)

func TestSliceIter(t *testing.T) {
	testdata := []int{1, 2, 3, 4, 5, 6, 7}

	iter := NewSliceIter(testdata)
	slice := testdata[:]
	for len(slice) != 0 {
		assertIterIsEmpty(t, false, iter)
		assertIterPeek(t, slice[0], iter)
		iter.Next()
		slice = slice[1:]
	}

	assertIterIsEmpty(t, true, &iter)
}

func assertIterIsEmpty(t *testing.T, expected bool, iter ConstIter) {
	t.Helper()
	if iter.IsEmpty() != expected {
		if expected {
			t.Errorf("Iter %#v is not empty", iter)
		} else {
			t.Errorf("Iter %#v is empty", iter)
		}
	}
}

func assertIterPeek(t *testing.T, expected int, iter ConstIter) {
	t.Helper()
	peeked := iter.Peek().(int)
	if peeked != expected {
		t.Errorf("Iter %#v: Peek returned %v, expected %v", iter, peeked, expected)
	}
}
