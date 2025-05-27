package condition

import (
	"strings"
)

type EqualsEvaluator struct {
	Value string
}

func NewEquals(value string) Evaluator {
	return &EqualsEvaluator{Value: value}
}

func (e *EqualsEvaluator) Evaluate(output string) bool {
	return strings.TrimSpace(output) == e.Value
}
