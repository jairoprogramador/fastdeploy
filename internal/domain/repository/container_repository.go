package repository

import (
	"deploy/internal/domain/variable"
)

type ContainerRepository interface {
	CreateFile(pathFile string, content string) error
	CreateDockerfile(pathFile, pathTemplate string, store *variable.VariableStore) error
	CreateDockerCompose(pathFile, pathTemplate string, store *variable.VariableStore) error
}
