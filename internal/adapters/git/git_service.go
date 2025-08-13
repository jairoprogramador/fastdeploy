package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/git"
)

type GitServiceImpl struct {
	pathResolver GitPathResolver
}

func NewGitService(
	pathResolver GitPathResolver,
) domain.GitService {
	return &GitServiceImpl{
		pathResolver: pathResolver,
	}
}

func (gs *GitServiceImpl) Clone(repositoryURL string) error {
	directoryPath, err := gs.pathResolver.GetDirectoryPath(repositoryURL)
	if err != nil {
		return err
	}

	gitDir := filepath.Join(directoryPath, ".git")

	if _, err := os.Stat(gitDir); err == nil {
		fmt.Printf("El repositorio ya está clonado en '%s'. Omitiendo la clonación.\n", directoryPath)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error al verificar el directorio .git: %w", err)
	}

	if err := os.MkdirAll(directoryPath, 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio '%s': %w", directoryPath, err)
	}

	cmd := exec.Command("git", "clone", repositoryURL, directoryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Clonando repositorio '%s' en '%s'...\n", repositoryURL, directoryPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al ejecutar git clone: %w", err)
	}

	return nil
}

func (gs *GitServiceImpl) IsCloned(repositoryURL string) bool {
	directoryPath, err := gs.pathResolver.GetDirectoryPath(repositoryURL)
	if err != nil {
		fmt.Printf("Error al obtener el directorio de configuración: %v\n", err)
		return false
	}

	gitDir := filepath.Join(directoryPath, ".git")

	if _, err := os.Stat(gitDir); err == nil {
		return true
	}
	return false
}
