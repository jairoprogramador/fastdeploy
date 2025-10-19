package application

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jairoprogramador/fastdeploy/internal/application/ports"
)

type UserInputProvider struct{}

func NewUserInputProvider() ports.UserInputProvider {
	return &UserInputProvider{}
}

func (p *UserInputProvider) Prompt(_ context.Context, question, defaultValue string) (string, error) {
	var response string
	prompt := &survey.Input{
		Message: question,
		Default: defaultValue,
	}
	err := survey.AskOne(prompt, &response)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (p *UserInputProvider) Confirm(_ context.Context, question string, defaultValue bool) (bool, error) {
	var response bool
	prompt := &survey.Confirm{
		Message: question,
		Default: defaultValue,
	}
	err := survey.AskOne(prompt, &response)
	return response, err
}
