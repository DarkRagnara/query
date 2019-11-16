package query

import (
	"testing"
)

type testdata struct {
	id   int
	name string
}

type testdb []testdata

func createData() testdb {
	return testdb{
		{id: 1, name: "abc"},
		{id: 2, name: "abc"},
		{id: 3, name: "ABC"},
		{id: 4, name: "def"},
		{id: 5, name: "abcdef"},
		{id: 6, name: "defabc"},
	}
}

type test struct {
	q         string
	resultIdx []int
}

func TestQueryWithEqualsInt(t *testing.T) {
	db := createData()
	proc := NewProcessor(db)

	tests := []test{
		{"id = '3'", []int{2}},
		{"id = '5'", []int{4}},
	}

	for _, test := range tests {
		query, err := proc.Build(test.q)

		if err != nil {
			t.Errorf("Error while building query {%v}: %v", test.q, err)
			continue
		}

		results := query.Run()
		assertResults(t, results, db, test)
	}
}

func assertResults(t *testing.T, results []interface{}, db testdb, test test) {
	t.Helper()
	expected := getExpected(db, test)

	if len(results) != len(expected) {
		t.Errorf("Expected %#v, but got %#v (for query {%v})", expected, results, test.q)
	}
}

func getExpected(db testdb, test test) []testdata {
	expected := []testdata{}
	for _, idx := range test.resultIdx {
		expected = append(expected, db[idx])
	}
	return expected
}
