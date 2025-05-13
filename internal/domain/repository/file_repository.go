package repository

import (
	"deploy/internal/domain/variable"
)

type FileRepository interface {
	GetFullPathDockerComposeTemplate(store *variable.VariableStore) string
	GetFullPathDockerCompose(store *variable.VariableStore) string
	GetFullPathDockerfileTemplate(store *variable.VariableStore) string
	GetFullPathDockerfile(store *variable.VariableStore) string
	ExistsFile(path string) bool
	DeleteFile(path string) error
}