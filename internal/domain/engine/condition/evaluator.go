package condition

type Evaluator interface {
	Evaluate(output string) bool
}
