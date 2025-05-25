package condition

import (
	"strings"
)

type ContainsEvaluator struct {
	Value string
}

func NewContainsEvaluator(value string) *ContainsEvaluator {
	return &ContainsEvaluator{Value: value}
}

func (e *ContainsEvaluator) Evaluate(output string) bool {
	return strings.Contains(output, e.Value)
}

