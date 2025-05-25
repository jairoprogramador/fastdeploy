package condition

import (
	"strings"
)

type NotEmptyEvaluator struct{}

func NewNotEmptyEvaluator() *NotEmptyEvaluator {
	return &NotEmptyEvaluator{}
}

func (e *NotEmptyEvaluator) Evaluate(output string) bool {
	return strings.TrimSpace(output) != ""
}
