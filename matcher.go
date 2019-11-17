package query

import (
	"strconv"
	"strings"
)

type Matcher interface {
	MatchEquals(interface{}, interface{}) bool
	TransformValue(string) (interface{}, error)
}

type IntEqualsMatcher struct{}

func (i IntEqualsMatcher) MatchEquals(a interface{}, b interface{}) bool {
	aInt := a.(int)
	bInt := b.(int)
	return aInt == bInt
}

func (i IntEqualsMatcher) TransformValue(s string) (interface{}, error) {
	return strconv.Atoi(s)
}

type StringEqualsMatcher struct{}

func (s StringEqualsMatcher) MatchEquals(a interface{}, b interface{}) bool {
	aStr := a.(string)
	bStr := b.(string)
	return strings.EqualFold(aStr, bStr)
}

func (s StringEqualsMatcher) TransformValue(str string) (interface{}, error) {
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
