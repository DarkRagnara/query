package query

type Processor struct {
}

func NewProcessor(db interface{}) Processor {
	return Processor{}
}

func (p Processor) Build(query string) (Query, error) {
	return Query{}, nil
}
