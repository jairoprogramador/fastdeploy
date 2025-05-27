package condition

import (
	"strings"
)

type EmptyEvaluator struct{}

func NewEmpty() Evaluator {
	return &EmptyEvaluator{}
}

func (e *EmptyEvaluator) Evaluate(output string) bool {
	return strings.TrimSpace(output) == ""
}
