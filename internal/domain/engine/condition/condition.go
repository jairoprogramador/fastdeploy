package condition

type ConditionType string

const (
	NotEmpty ConditionType = "not_empty"
	Empty    ConditionType = "empty"
	Equals   ConditionType = "equals"
	Contains ConditionType = "contains"
	Matches  ConditionType = "matches"
)

type ConditionEvaluator interface {
	Evaluate(output string) bool
}
