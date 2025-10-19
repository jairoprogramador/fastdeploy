package vos

type DecisionType int

const (
	ShouldExecute DecisionType = iota
	ShouldSkip
)

type ExecutionDecision struct {
	decisionType DecisionType
	reason       string
}

func Execute(reason string) ExecutionDecision {
	return ExecutionDecision{decisionType: ShouldExecute, reason: reason}
}

func Skip(reason string) ExecutionDecision {
	return ExecutionDecision{decisionType: ShouldSkip, reason: reason}
}

func (d ExecutionDecision) ShouldExecute() bool {
	return d.decisionType == ShouldExecute
}

func (d ExecutionDecision) Reason() string {
	return d.reason
}