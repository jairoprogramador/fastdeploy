package entities

import (
	"github.com/stretchr/testify/mock"

	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
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

func (m *MockVariableResolver) ProcessTemplateFile(srcPath, destPath string, variables map[string]vos.Variable) error {
	args := m.Called(srcPath, destPath, variables)
	return args.Error(0)
}