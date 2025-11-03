package vos

type ExecutionStatus int

const (
	ExecutionStatusInProgress ExecutionStatus = iota
	ExecutionStatusSuccessful
	ExecutionStatusFailed
)

func (s ExecutionStatus) String() string {
	switch s {
	case ExecutionStatusInProgress:
		return StatusInProgress.String()
	case ExecutionStatusSuccessful:
		return StatusSuccessful.String()
	case ExecutionStatusFailed:
		return StatusFailed.String()
	default:
		return StatusUnknown.String()
	}
}

func ExecutionStatusFromString(status string) ExecutionStatus {
	switch status {
	case StatusInProgress.String():
		return ExecutionStatusInProgress
	case StatusSuccessful.String():
		return ExecutionStatusSuccessful
	case StatusFailed.String():
		return ExecutionStatusFailed
	default:
		return ExecutionStatus(99)
	}
}
