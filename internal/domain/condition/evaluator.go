package condition

import (
	"regexp"
	"strings"
	"sync"
)

// Evaluadores singleton
var (
	notEmptyEvaluator *NotEmptyEvaluator
	emptyEvaluator    *EmptyEvaluator
	onceNotEmpty      sync.Once
	onceEmpty         sync.Once
)

type NotEmptyEvaluator struct{}
type EmptyEvaluator struct{}
type EqualsEvaluator struct {
	Value string
}
type ContainsEvaluator struct {
	Value string
}
type MatchesEvaluator struct {
	Pattern *regexp.Regexp
}

// GetNotEmptyEvaluator retorna la instancia única del evaluador
func GetNotEmptyEvaluator() *NotEmptyEvaluator {
	onceNotEmpty.Do(func() {
		notEmptyEvaluator = &NotEmptyEvaluator{}
	})
	return notEmptyEvaluator
}

// GetEmptyEvaluator retorna la instancia única del evaluador
func GetEmptyEvaluator() *EmptyEvaluator {
	onceEmpty.Do(func() {
		emptyEvaluator = &EmptyEvaluator{}
	})
	return emptyEvaluator
}

func (e *NotEmptyEvaluator) Evaluate(output string) (bool, error) {
	return strings.TrimSpace(output) != "", nil
}

func (e *EmptyEvaluator) Evaluate(output string) (bool, error) {
	return strings.TrimSpace(output) == "", nil
}

func (e *EqualsEvaluator) Evaluate(output string) (bool, error) {
	return strings.TrimSpace(output) == e.Value, nil
}

func (e *ContainsEvaluator) Evaluate(output string) (bool, error) {
	return strings.Contains(output, e.Value), nil
}

func (e *MatchesEvaluator) Evaluate(output string) (bool, error) {
	return e.Pattern.MatchString(strings.TrimSpace(output)), nil
}
