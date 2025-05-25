package condition

import (
	"regexp"
	"strings"
)

type MatchesEvaluator struct {
	Pattern *regexp.Regexp
}

func NewMatchesEvaluator(pattern *regexp.Regexp) *MatchesEvaluator {
	return &MatchesEvaluator{Pattern: pattern}
}

func (e *MatchesEvaluator) Evaluate(output string) bool {
	return e.Pattern.MatchString(strings.TrimSpace(output))
}
