package condition

import (
	"regexp"
	"strings"
)

type NotEmptyEvaluator struct{}
type EmptyEvaluator struct{}
type EqualsEvaluator struct {
	Value string
}
type ContainsEvaluator struct {
	Value string
}
type MatchesEvaluator struct {
	Pattern *regexp.Regexp
}

func NewNotEmptyEvaluator() *NotEmptyEvaluator {
	return &NotEmptyEvaluator{}
}

func NewEmptyEvaluator() *EmptyEvaluator {
	return &EmptyEvaluator{}
}

func (e *NotEmptyEvaluator) Evaluate(output string) (bool, error) {
	return strings.TrimSpace(output) != "", nil
}

func (e *EmptyEvaluator) Evaluate(output string) (bool, error) {
	return strings.TrimSpace(output) == "", nil
}

func (e *EqualsEvaluator) Evaluate(output string) (bool, error) {
	return strings.TrimSpace(output) == e.Value, nil
}

func (e *ContainsEvaluator) Evaluate(output string) (bool, error) {
	return strings.Contains(output, e.Value), nil
}

func (e *MatchesEvaluator) Evaluate(output string) (bool, error) {
	return e.Pattern.MatchString(strings.TrimSpace(output)), nil
}
