package condition

import (
	"regexp"
	"strings"
)

type ConditionFactory struct{}

func NewConditionFactory() *ConditionFactory {
	return &ConditionFactory{}
}

func (f *ConditionFactory) CreateEvaluator(conditionStr string) ConditionEvaluator {
	parts := strings.SplitN(conditionStr, ":", 2)

	switch parts[0] {
		case string(NotEmpty):
			return NewNotEmptyEvaluator()
		case string(Empty):
			return NewEmptyEvaluator()
		case string(Equals):
			return &EqualsEvaluator{Value: strings.TrimSpace(parts[1])}
		case string(Contains):
			return &ContainsEvaluator{Value: strings.TrimSpace(parts[1])}
		case string(Matches):
			pattern, _ := regexp.Compile(strings.TrimSpace(parts[1]))
			return &MatchesEvaluator{Pattern: pattern}
		default:
			return nil
	}
}
