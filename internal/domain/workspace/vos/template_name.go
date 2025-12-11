package vos

import "errors"

type TemplateName struct {
	value string
}

func NewTemplateName(value string) (TemplateName, error) {
	if value == "" {
		return TemplateName{}, errors.New("template name cannot be empty")
	}
	return TemplateName{value: value}, nil
}

func (n TemplateName) String() string {
	return n.value
}
