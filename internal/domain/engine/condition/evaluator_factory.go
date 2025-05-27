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
		return NewNotEmpty()
	case string(Empty):
		return NewEmpty()
	case string(Equals):
		return NewEquals(strings.TrimSpace(parts[1]))
	case string(Contains):
		return NewContains(strings.TrimSpace(parts[1]))
	case string(Matches):
		pattern, _ := regexp.Compile(strings.TrimSpace(parts[1]))
		return NewMatches(pattern)
	default:
		return nil
	}
}
