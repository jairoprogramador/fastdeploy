package console

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jairoprogramador/fastdeploy/newinternal/application/ports"
)

// UserInputProvider implementa la interfaz ports.UserInputProvider.
type UserInputProvider struct{}

// NewUserInputProvider crea una nueva instancia.
func NewUserInputProvider() ports.UserInputProvider {
	return &UserInputProvider{}
}

// Prompt utiliza la librería 'survey' para presentar una pregunta al usuario.
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

// Confirm presenta una pregunta de sí/no.
func (p *UserInputProvider) Confirm(_ context.Context, question string, defaultValue bool) (bool, error) {
	var response bool
	prompt := &survey.Confirm{
		Message: question,
		Default: defaultValue,
	}
	err := survey.AskOne(prompt, &response)
	return response, err
}
