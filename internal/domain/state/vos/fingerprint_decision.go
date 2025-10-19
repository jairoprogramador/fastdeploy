package vos

type DecisionType int

const (
	ShouldExecute DecisionType = iota
	ShouldSkip
)

type FingerprintDecision struct {
	decisionType DecisionType
	reason       string
}

func Execute(reason string) FingerprintDecision {
	return FingerprintDecision{decisionType: ShouldExecute, reason: reason}
}

func Skip(reason string) FingerprintDecision {
	return FingerprintDecision{decisionType: ShouldSkip, reason: reason}
}

func (d FingerprintDecision) ShouldExecute() bool {
	return d.decisionType == ShouldExecute
}

func (d FingerprintDecision) Reason() string {
	return d.reason
}
