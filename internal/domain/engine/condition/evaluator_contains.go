package condition

import (
	"strings"
)

type ContainsEvaluator struct {
	Value string
}

func NewContains(value string) Evaluator {
	return &ContainsEvaluator{Value: value}
}

func (e *ContainsEvaluator) Evaluate(output string) bool {
	return strings.Contains(output, e.Value)
}
