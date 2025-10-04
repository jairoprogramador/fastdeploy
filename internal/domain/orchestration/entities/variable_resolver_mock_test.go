package entities

import (
	"github.com/stretchr/testify/mock"

	deploymentvos "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
)

// MockVariableResolver es un mock para la interfaz services.VariableResolver.
// Nos permite simular su comportamiento en los tests.
type MockVariableResolver struct {
	mock.Mock
}

func (m *MockVariableResolver) ExtractVariable(probe deploymentvos.OutputProbe, text string) (vos.Variable, bool, error) {
	args := m.Called(probe, text)
	return args.Get(0).(vos.Variable), args.Bool(1), args.Error(2)
}

func (m *MockVariableResolver) Interpolate(template string, variables map[string]vos.Variable) (string, error) {
	args := m.Called(template, variables)
	return args.String(0), args.Error(1)
}

func (m *MockVariableResolver) ProcessTemplate(pathFile string, variables map[string]vos.Variable) error {
	args := m.Called(pathFile, variables)
	return args.Error(0)
}