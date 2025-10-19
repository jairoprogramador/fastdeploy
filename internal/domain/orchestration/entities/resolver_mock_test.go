package entities

import (
	"github.com/stretchr/testify/mock"

	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
)

type MockResolver struct {
	mock.Mock
}

func (m *MockResolver) ResolveOutput(probe vos.Output, text string) (vos.Output, bool, error) {
	args := m.Called(probe, text)
	return args.Get(0).(vos.Output), args.Bool(1), args.Error(2)
}

func (m *MockResolver) ResolveTemplate(template string, variables map[string]vos.Output) (string, error) {
	args := m.Called(template, variables)
	return args.String(0), args.Error(1)
}

func (m *MockResolver) ResolvePath(path string, variables map[string]vos.Output) error {
	args := m.Called(path, variables)
	return args.Error(0)
}
