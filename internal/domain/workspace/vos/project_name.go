package vos

import "errors"

type ProjectName struct {
	value string
}

func NewProjectName(value string) (ProjectName, error) {
	if value == "" {
		return ProjectName{}, errors.New("project name cannot be empty")
	}
	return ProjectName{value: value}, nil
}

func (n ProjectName) String() string {
	return n.value
}
