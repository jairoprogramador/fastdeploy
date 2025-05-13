package repository

import (
	"deploy/internal/domain/repository"
	"deploy/internal/infrastructure/tools"
	"sync"
)

type variableRepositoryImpl struct{}

var (
	instanceVariableRepository     repository.VariableRepository
	instanceOnceVariableRepository sync.Once
)

func GetVariableRepository() repository.VariableRepository {
	instanceOnceVariableRepository.Do(func() {
		instanceVariableRepository = &variableRepositoryImpl{}
	})
	return instanceVariableRepository
}

func (s *variableRepositoryImpl) GetCommitHash() (string, error) {
	return tools.GetCommitHash()
}

func (s *variableRepositoryImpl) GetCommitAuthor(commitHash string) (string, error) {
	return tools.GetCommitAuthor(commitHash)
}

func (s *variableRepositoryImpl) GetCommitMessage(commitHash string) (string, error) {
	return tools.GetCommitMessage(commitHash)
}
