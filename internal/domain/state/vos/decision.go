package vos

type DecisionType int

const (
	ShouldExecute DecisionType = iota
	ShouldSkip
)

type Decision struct {
	decisionType DecisionType
	reason       string
}

func Execute(reason string) Decision {
	return Decision{decisionType: ShouldExecute, reason: reason}
}

func Skip(reason string) Decision {
	return Decision{decisionType: ShouldSkip, reason: reason}
}

func (d Decision) ShouldExecute() bool {
	return d.decisionType == ShouldExecute
}

func (d Decision) Reason() string {
	return d.reason
}
