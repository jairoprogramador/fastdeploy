package condition

import (
	"regexp"
	"strings"
)

type EvaluatorFactory struct{}

func NewEvaluatorFactory() *EvaluatorFactory {
	return &EvaluatorFactory{}
}

func (f *EvaluatorFactory) CreateEvaluator(conditionStr string) Evaluator {
	parts := strings.SplitN(conditionStr, ":", 2)

	switch parts[0] {
	case string(NotEmpty):
		return NewNotEmptyEvaluator()
	case string(Empty):
		return NewEmptyEvaluator()
	case string(Equals):
		return NewEqualsEvaluator(strings.TrimSpace(parts[1]))
	case string(Contains):
		return NewContainsEvaluator(strings.TrimSpace(parts[1]))
	case string(Matches):
		pattern, _ := regexp.Compile(strings.TrimSpace(parts[1]))
		return NewMatchesEvaluator(pattern)
	default:
		return nil
	}
}
