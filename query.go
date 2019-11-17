package query

type Query struct {
	matcher func(interface{}) bool
	db      Queryable
}

func (q Query) Run() []interface{} {
	results := []interface{}{}
	iter := q.db.Iter()
	for !iter.IsEmpty() {
		front := iter.Peek()
		if q.matcher(front) {
			results = append(results, front)
		}
		iter.Next()
	}
	return results
}
