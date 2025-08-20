package project

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"os/exec"

	"github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
)

type Git struct{}

func NewGit() ports.Git {
	return &Git{}
}

func (g *Git) Clone(url string, nameRepository string) error {
	directoryPath, err := g.getDirectoryPath(nameRepository)
	if err != nil {
		return fmt.Errorf("clone repository failed, get directory path error: %w", err)
	}

	directoryGit := filepath.Join(directoryPath, ".git")

	if _, err := os.Stat(directoryGit); err == nil {
		fmt.Printf("El repositorio ya está clonado en '%s'. Omitiendo la clonación.\n", directoryPath)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("clone repository failed, error al verificar el directorio .git: %w", err)
	}

	if err := os.MkdirAll(directoryPath, 0755); err != nil {
		return fmt.Errorf("clone repository failed, no se pudo crear el directorio '%s': %w", directoryPath, err)
	}

	cmd := exec.Command("git", "clone", url, directoryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Clonando repositorio '%s' en '%s'...\n", url, directoryPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clone repository failed, error al ejecutar git clone: %w", err)
	}

	return nil
}

func (g *Git) IsCloned(nameRepository string) (bool, error) {
	directoryPath, err := g.getDirectoryPath(nameRepository)
	if err != nil {
		return false, fmt.Errorf("is cloned repository failed, get directory path error: %w", err)
	}

	directoryGit := filepath.Join(directoryPath, ".git")

	if _, err := os.Stat(directoryGit); err == nil {
		return true, nil
	}
	return false, nil
}

func (g *Git) getDirectoryPath(nameRepository string) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}

	directoryPath := filepath.Join(currentUser.HomeDir, constants.FastDeployDir)

	return filepath.Join(directoryPath, nameRepository), nil
}
