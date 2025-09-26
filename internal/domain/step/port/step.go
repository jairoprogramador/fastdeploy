package port

import "github.com/jairoprogramador/fastdeploy/internal/domain/command/values"


type StepPort interface {
	LoadCommands(pathStepFile string) ([]values.CommandValue, error)
}