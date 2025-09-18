package service

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"os/exec"

	"github.com/jairoprogramador/fastdeploy/internal/domain/project/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
)

type GitManager struct{}

func NewGitManager() port.GitManager {
	return &GitManager{}
}

func (g *GitManager) Clone(url string, nameRepository string) error {
	directoryPath, err := g.getDirectoryPath(nameRepository)
	if err != nil {
		return err
	}

	directoryGit := filepath.Join(directoryPath, ".git")

	if _, err := os.Stat(directoryGit); err == nil {
		fmt.Printf("El repositorio de despliegue ya está clonado en '%s'. Omitiendo la clonación.\n", directoryPath)
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(directoryPath, 0755); err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", url, directoryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Clonando repositorio '%s' en '%s'...\n", url, directoryPath)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (g *GitManager) IsCloned(nameRepository string) (bool, error) {
	directoryPath, err := g.getDirectoryPath(nameRepository)
	if err != nil {
		return false, err
	}

	directoryGit := filepath.Join(directoryPath, ".git")

	if _, err := os.Stat(directoryGit); err == nil {
		return true, nil
	}
	return false, nil
}

func (g *GitManager) getDirectoryPath(nameRepository string) (string, error) {
	directoryPath, err := g.getHomeDirPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(directoryPath, nameRepository), nil
}

func (g *GitManager) getHomeDirPath() (string, error) {
	if fastDeployHome := os.Getenv("FASTDEPLOY_HOME"); fastDeployHome != "" {
		return fastDeployHome, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}
	return filepath.Join(currentUser.HomeDir, constants.FastDeployDir), nil
}