package condition

import (
	"strings"
)

type EqualsEvaluator struct {
	Value string
}

func NewEqualsEvaluator(value string) *EqualsEvaluator {
	return &EqualsEvaluator{Value: value}
}

func (e *EqualsEvaluator) Evaluate(output string) bool {
	return strings.TrimSpace(output) == e.Value
}

