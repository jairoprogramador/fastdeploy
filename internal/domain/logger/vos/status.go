package vos

import (
	"fmt"
)

type Status int

const (
	Pending Status = iota
	Running
	Success
	Failure
	Skipped
	Cached
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Failure:
		return "Failure"
	case Skipped:
		return "Skipped"
	case Cached:
		return "Cached"
	default:
		return "Unknown"
	}
}

func NewStatusFromString(status string) (Status, error) {
	switch status {
	case "Pending":
		return Pending, nil
	case "Running":
		return Running, nil
	case "Success":
		return Success, nil
	case "Failure":
		return Failure, nil
	case "Skipped":
		return Skipped, nil
	case "Cached":
		return Cached, nil
	default:
		return 0, fmt.Errorf("invalid status: %s", status)
	}
}
