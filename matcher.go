package query

import (
	"strconv"
	"strings"
)

type Matcher interface {
	MatchEquals(interface{}, interface{}) bool
	TransformValue(string) (interface{}, error)
}

type IntMatcher struct{}

func (i IntMatcher) MatchEquals(a interface{}, b interface{}) bool {
	aInt := a.(int)
	bInt := b.(int)
	return aInt == bInt
}

func (i IntMatcher) TransformValue(s string) (interface{}, error) {
	return strconv.Atoi(s)
}

type StringMatcher struct{}

func (s StringMatcher) MatchEquals(a interface{}, b interface{}) bool {
	aStr := a.(string)
	bStr := b.(string)
	return strings.EqualFold(aStr, bStr)
}

func (s StringMatcher) TransformValue(str string) (interface{}, error) {
	return str, nil
}

type matcherFunc = func(interface{}, interface{}) bool

type operator int

const (
	opEquals operator = iota
)

func funcByOP(m Matcher, op operator) matcherFunc {
	switch op {
	case opEquals:
		return m.MatchEquals
	}
	panic("Unknown Operator")
}
