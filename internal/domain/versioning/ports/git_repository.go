package ports

import (
	"context"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/versioning/vos"
)

// GitRepository define la interfaz para interactuar con un repositorio Git.
type GitRepository interface {
	// GetLastCommit obtiene el último commit de la rama actual.
	GetLastCommit(ctx context.Context, repoPath string) (*vos.Commit, error)

	// GetCommitsSinceTag obtiene la lista de commits desde el último tag semántico.
	GetCommitsSinceTag(ctx context.Context, repoPath string, lastTag string) ([]*vos.Commit, error)

	// GetLastSemverTag obtiene el último tag que sigue el formato de versionado semántico.
	GetLastSemverTag(ctx context.Context, repoPath string) (string, error)
}
