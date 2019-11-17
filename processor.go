package query

import (
	"bitbucket.org/ragnara/pars/v2"
	"unicode"
)

type Processor struct {
	db     Queryable
	fields []Field
}

func NewProcessor(db Queryable) Processor {
	fields := db.Fields()
	return Processor{db: db, fields: fields}
}

func (p Processor) Build(query string) (*Query, error) {
	matcher, err := pars.ParseString(query, p.parser())
	if err != nil {
		return nil, err
	}

	return &Query{db: p.db, matcher: matcher.(func(interface{}) bool)}, nil
}

func (p Processor) parser() pars.Parser {
	return pars.Dispatch(p.termClause())
}

func (p Processor) termClause() pars.DispatchClause {
	return pars.DescribeClause{
		DispatchClause: pars.Clause{p.fieldParser()},
		Description:    "query term"}
}

func (p Processor) fieldParser() pars.Parser {
	clauses := make([]pars.DispatchClause, 0, len(p.fields))
	for _, field := range p.fields {
		clauses = append(clauses, field)
	}
	return pars.Dispatch(clauses...)
}

type Queryable interface {
	Fields() []Field
	Iter() Iter
}

type Field struct {
	Name    string
	Getter  func(interface{}) interface{}
	Matcher Matcher
}

func (f Field) Parsers() []pars.Parser {
	return []pars.Parser{
		identifierParser(f.Name),
		opParser(f.Matcher),
		valueParser(f.Matcher),
	}
}

func (f Field) TransformResult(v []interface{}) interface{} {
	fn := v[1].(matcherFunc)
	compareTo := v[2]
	return func(val interface{}) bool {
		return fn(f.Getter(val), compareTo)
	}
}

func (f Field) TransformError(err error) error {
	return err
}

func opParser(m Matcher) pars.Parser {
	return pars.Or(
		pars.Transformer(pars.Char('='), func(interface{}) (interface{}, error) {
			return funcByOP(m, opEquals), nil
		}),
		pars.Transformer(pars.String("<>"), func(interface{}) (interface{}, error) {
			return funcByOP(m, opNotEquals), nil
		}),
		pars.Transformer(pars.Char('<'), func(interface{}) (interface{}, error) {
			return funcByOP(m, opLessThan), nil
		}),
		pars.Transformer(pars.Char('>'), func(interface{}) (interface{}, error) {
			return funcByOP(m, opGreaterThan), nil
		}),
		pars.Transformer(likeKeywordParser(), func(interface{}) (interface{}, error) {
			return funcByOP(m, opLike), nil
		}))
}

func wholeWordParser(p pars.Parser) pars.Parser {
	return pars.SwallowWhitespace(
		pars.Except(
			p.Clone(),
			pars.Seq(
				p.Clone(),
				pars.CharPred(func(r rune) bool {
					return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
				}))))
}

func likeKeywordParser() pars.Parser {
	return wholeWordParser(pars.StringCI("like"))
}

func identifierParser(name string) pars.Parser {
	return wholeWordParser(pars.String(name))
}

func valueParser(m Matcher) pars.Parser {
	return pars.Transformer(
		pars.SwallowWhitespace(pars.DelimitedString("'", "'")),
		func(v interface{}) (interface{}, error) {
			return m.TransformValue(v.(string))
		})
}
