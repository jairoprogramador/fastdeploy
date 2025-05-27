package condition

import (
	"regexp"
	"strings"
)

type MatchesEvaluator struct {
	Pattern *regexp.Regexp
}

func NewMatches(pattern *regexp.Regexp) Evaluator {
	return &MatchesEvaluator{Pattern: pattern}
}

func (e *MatchesEvaluator) Evaluate(output string) bool {
	return e.Pattern.MatchString(strings.TrimSpace(output))
}
