package query

import (
	"testing"
)

type testdata struct {
	id   int
	id2  int
	name string
}

type testdb []testdata

func (db testdb) Fields() []Field {
	return []Field{
		{"id", func(entry interface{}) interface{} {
			return entry.(testdata).id
		}, IntMatcher{}},
		{"id2", func(entry interface{}) interface{} {
			return entry.(testdata).id2
		}, IntMatcher{}},
		{"name", func(entry interface{}) interface{} {
			return entry.(testdata).name
		}, StringMatcher{}},
	}
}

func (db testdb) Iter() Iter {
	return &SliceIter{db}
}

func createData() testdb {
	return testdb{
		{id: 1, id2: 6, name: "abc"},
		{id: 2, id2: 5, name: "abc"},
		{id: 3, id2: 4, name: "ABC"},
		{id: 4, id2: 3, name: "def"},
		{id: 5, id2: 2, name: "abcdef"},
		{id: 6, id2: 1, name: "defabc"},
		{id: 7, id2: 0, name: "defabcdef"},
	}
}

type test struct {
	q         string
	resultIdx []int
}

func TestQueryWithLikeString(t *testing.T) {
	tests := []test{
		{"name like '%abc'", []int{0, 1, 2, 5}},
		{"name LIKE 'abc%'", []int{0, 1, 2, 4}},
		{"name like '%abc%'", []int{0, 1, 2, 4, 5, 6}},
		{"name LiKe '%'", []int{0, 1, 2, 3, 4, 5, 6}},
	}

	runTests(t, tests)
}

func TestQueryWithEqualsInt(t *testing.T) {
	tests := []test{
		{"id = '3'", []int{2}},
		{"id = '5'", []int{4}},
		{"id2 = '3'", []int{3}},
		{"id2 = '5'", []int{1}},
	}

	runTests(t, tests)
}

func TestQueryWithEqualsString(t *testing.T) {
	tests := []test{
		{"name = 'abc'", []int{0, 1, 2}},
		{"name = 'def'", []int{3}},
	}

	runTests(t, tests)
}

func TestQueryWithLessThanInt(t *testing.T) {
	tests := []test{
		{"id < '3'", []int{0, 1}},
		{"id < '5'", []int{0, 1, 2, 3}},
		{"id2 < '3'", []int{4, 5, 6}},
		{"id2 < '5'", []int{2, 3, 4, 5, 6}},
	}

	runTests(t, tests)
}

func TestQueryWithLessThanString(t *testing.T) {
	tests := []test{
		{"name < 'abc'", []int{}},
		{"name < 'def'", []int{0, 1, 2, 4}},
		{"name < 'abd'", []int{0, 1, 2, 4}},
	}

	runTests(t, tests)
}

func TestQueryWithLessThanOrEqualsInt(t *testing.T) {
	tests := []test{
		{"id <= '3'", []int{0, 1, 2}},
		{"id <= '5'", []int{0, 1, 2, 3, 4}},
		{"id2 <= '3'", []int{3, 4, 5, 6}},
		{"id2 <= '5'", []int{1, 2, 3, 4, 5, 6}},
	}

	runTests(t, tests)
}

func TestQueryWithLessThanOrEqualsString(t *testing.T) {
	tests := []test{
		{"name <= 'abc'", []int{0, 1, 2}},
		{"name <= 'def'", []int{0, 1, 2, 3, 4}},
		{"name <= 'abd'", []int{0, 1, 2, 4}},
	}

	runTests(t, tests)
}

func TestQueryWithGreaterThanInt(t *testing.T) {
	tests := []test{
		{"id > '3'", []int{3, 4, 5, 6}},
		{"id > '5'", []int{5, 6}},
		{"id2 > '3'", []int{0, 1, 2}},
		{"id2 > '5'", []int{0}},
	}

	runTests(t, tests)
}

func TestQueryWithGreaterThanString(t *testing.T) {
	tests := []test{
		{"name > 'abc'", []int{3, 4, 5, 6}},
		{"name > 'def'", []int{5, 6}},
		{"name > 'abd'", []int{3, 5, 6}},
	}

	runTests(t, tests)
}

func TestQueryWithGreaterThanOrEqualsInt(t *testing.T) {
	tests := []test{
		{"id >= '3'", []int{2, 3, 4, 5, 6}},
		{"id >= '5'", []int{4, 5, 6}},
		{"id2 >= '3'", []int{0, 1, 2, 3}},
		{"id2 >= '5'", []int{0, 1}},
	}

	runTests(t, tests)
}

func TestQueryWithGreaterThanOrEqualsString(t *testing.T) {
	tests := []test{
		{"name >= 'abc'", []int{0, 1, 2, 3, 4, 5, 6}},
		{"name >= 'def'", []int{3, 5, 6}},
		{"name >= 'abd'", []int{3, 5, 6}},
	}

	runTests(t, tests)
}

func TestQueryWithNotEqualsInt(t *testing.T) {
	tests := []test{
		{"id <> '3'", []int{0, 1, 3, 4, 5, 6}},
		{"id <> '5'", []int{0, 1, 2, 3, 5, 6}},
		{"id2 <> '3'", []int{0, 1, 2, 4, 5, 6}},
		{"id2 <> '5'", []int{0, 2, 3, 4, 5, 6}},
	}

	runTests(t, tests)
}

func TestQueryWithNotEqualsString(t *testing.T) {
	tests := []test{
		{"name <> 'abc'", []int{3, 4, 5, 6}},
		{"name <> 'def'", []int{0, 1, 2, 4, 5, 6}},
	}

	runTests(t, tests)
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		q   string
		err string
	}{
		{"", "Identifier expected"},
		{"id something...", "Operator expected"},
		{"id = '3'garbage", "Expected EOF: Found byte 0x67"},
		{"id = '3", "Could not parse expected string \"'\": EOF"},
	}

	db := createData()
	proc := NewProcessor(db)
	for _, test := range tests {
		_, err := proc.Build(test.q)
		if err.Error() != test.err {
			t.Errorf("Query {%v}: Expected error %v, but got %v", test.q, test.err, err.Error())
		}
	}
}

func runTests(t *testing.T, tests []test) {
	t.Helper()
	db := createData()
	proc := NewProcessor(db)

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
