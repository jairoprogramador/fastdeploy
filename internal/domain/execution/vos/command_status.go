package vos

type CommandStatus int

const (
	CommandStatusPending CommandStatus = iota
	CommandStatusSuccessful
	CommandStatusFailed
)

func (s CommandStatus) String() string {
	switch s {
	case CommandStatusPending:
		return StatusPending.String()
	case CommandStatusSuccessful:
		return StatusSuccessful.String()
	case CommandStatusFailed:
		return StatusFailed.String()
	default:
		return StatusUnknown.String()
	}
}

func CommandStatusFromString(status string) CommandStatus {
	switch status {
	case StatusPending.String():
		return CommandStatusPending
	case StatusSuccessful.String():
		return CommandStatusSuccessful
	case StatusFailed.String():
		return CommandStatusFailed
	default:
		return CommandStatus(99)
	}
}