package condition

import (
	"fmt"
	"regexp"
	"strings"
)

type ConditionFactory struct{}

func NewConditionFactory() *ConditionFactory {
	return &ConditionFactory{}
}

func (f *ConditionFactory) CreateEvaluator(conditionStr string, output string) (ConditionEvaluator, error) {
	parts := strings.SplitN(conditionStr, ":", 2)

	switch parts[0] {
		case string(NotEmpty):
		return NewNotEmptyEvaluator(), nil
		case string(Empty):
		return NewEmptyEvaluator(), nil
		case string(Equals):
			if output == "" {
				return nil, fmt.Errorf("value is required for equals condition")
			}
			if parts[1] == "" {
				return nil, fmt.Errorf("value is required for equals condition")
			}
			return &EqualsEvaluator{Value: strings.TrimSpace(parts[1])}, nil
		case string(Contains):
			if output == "" {
				return nil, fmt.Errorf("value is required for contains condition")
			}
			if parts[1] == "" {
				return nil, fmt.Errorf("value is required for contains condition")
			}
			return &ContainsEvaluator{Value: strings.TrimSpace(parts[1])}, nil
		case string(Matches):
			if output == "" {
				return nil, fmt.Errorf("pattern is required for matches condition")
			}
			if parts[1] == "" {
				return nil, fmt.Errorf("value is required for matches condition")
			}
			pattern, err := regexp.Compile(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid regex pattern: %v", err)
			}
			return &MatchesEvaluator{Pattern: pattern}, nil
		default:
			return nil, fmt.Errorf("unknown condition type: %s", conditionStr)
	}
}
