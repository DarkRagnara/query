package query

import (
	"bitbucket.org/ragnara/pars/v2"
	"bitbucket.org/ragnara/query/iter"
	"errors"
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
	matcher, err := pars.ParseString(query, pars.DiscardRight(p.parser(), pars.EOF))
	if err != nil {
		return nil, err
	}

	return &Query{db: p.db, matcher: matcher.(queryFunc)}, nil
}

func (p Processor) parser() pars.Parser {
	return pars.Dispatch(
		orClause{&p},
		andClause{&p},
		p.termClause())
}

func (p Processor) termClause() pars.DispatchClause {
	return pars.DescribeClause{
		DispatchClause: pars.Clause{p.fieldParser()},
		Description:    "query term"}
}

func (p Processor) fieldParser() pars.Parser {
	clauses := make([]pars.DispatchClause, 0, len(p.fields)*2)
	for i := range p.fields {
		clauses = append(clauses, negated(p.fields[i]), p.fields[i])
	}
	clauses = append(clauses, pars.Clause{pars.Error(errors.New("Identifier expected"))})
	return pars.Dispatch(clauses...)
}

type Queryable interface {
	Fields() []Field
	Iter() iter.Iter
}

type orClause struct {
	p *Processor
}

func (o orClause) Parsers() []pars.Parser {
	return []pars.Parser{
		pars.DiscardRight(
			pars.Dispatch(
				andClause{o.p},
				o.p.termClause()),
			wholeWordParser(pars.StringCI("or"))),
		pars.Recursive(o.p.parser),
	}
}

func (o orClause) TransformResult(fns []interface{}) interface{} {
	left := fns[0].(queryFunc)
	right := fns[1].(queryFunc)
	return func(val interface{}) bool {
		return left(val) || right(val)
	}
}

func (o orClause) TransformError(err error) error {
	return err
}

type andClause struct {
	p *Processor
}

func (o andClause) Parsers() []pars.Parser {
	return []pars.Parser{
		pars.DiscardRight(o.p.fieldParser(), wholeWordParser(pars.StringCI("and"))),
		pars.Recursive(func() pars.Parser {
			return pars.Dispatch(andClause{o.p}, o.p.termClause())
		}),
	}
}

func (o andClause) TransformResult(fns []interface{}) interface{} {
	left := fns[0].(queryFunc)
	right := fns[1].(queryFunc)
	return func(val interface{}) bool {
		return left(val) && right(val)
	}
}

func (o andClause) TransformError(err error) error {
	return err
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
		pars.Transformer(pars.String("<="), func(interface{}) (interface{}, error) {
			return funcByOP(m, opLessThanOrEquals), nil
		}),
		pars.Transformer(pars.String(">="), func(interface{}) (interface{}, error) {
			return funcByOP(m, opGreaterThanOrEquals), nil
		}),
		pars.Transformer(pars.Char('<'), func(interface{}) (interface{}, error) {
			return funcByOP(m, opLessThan), nil
		}),
		pars.Transformer(pars.Char('>'), func(interface{}) (interface{}, error) {
			return funcByOP(m, opGreaterThan), nil
		}),
		pars.Transformer(likeKeywordParser(), func(interface{}) (interface{}, error) {
			return funcByOP(m, opLike), nil
		}),
		pars.Error(errors.New("Operator expected")))
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

type negatedField struct {
	Field
}

func (n negatedField) Parsers() []pars.Parser {
	return []pars.Parser{
		pars.DiscardLeft(wholeWordParser(pars.StringCI("not")), identifierParser(n.Field.Name)),
		opParser(n.Field.Matcher),
		valueParser(n.Field.Matcher),
	}
}

func (n negatedField) TransformResult(v []interface{}) interface{} {
	return func(val interface{}) bool {
		return !n.Field.TransformResult(v).(queryFunc)(val)
	}
}

func negated(f Field) negatedField {
	return negatedField{f}
}
