package vos

type StepStatus int

const (
	StepStatusPending StepStatus = iota
	StepStatusInProgress
	StepStatusSkipped
	StepStatusCached
	StepStatusSuccessful
	StepStatusFailed
)

func (s StepStatus) String() string {
	switch s {
	case StepStatusPending:
		return StatusPending.String()
	case StepStatusInProgress:
		return StatusInProgress.String()
	case StepStatusSkipped:
		return StatusSkipped.String()
	case StepStatusCached:
		return StatusCached.String()
	case StepStatusSuccessful:
		return StatusSuccessful.String()
	case StepStatusFailed:
		return StatusFailed.String()
	default:
		return StatusUnknown.String()
	}
}

func StepStatusFromString(status string) StepStatus {
	switch status {
	case StatusPending.String():
		return StepStatusPending
	case StatusInProgress.String():
		return StepStatusInProgress
	case StatusSkipped.String():
		return StepStatusSkipped
	case StatusCached.String():
		return StepStatusCached
	case StatusSuccessful.String():
		return StepStatusSuccessful
	case StatusFailed.String():
		return StepStatusFailed
	default:
		return StepStatus(99)
	}
}
