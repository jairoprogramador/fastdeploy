package repository

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
)

func CloneRepository(repositoryURL string) error {
	configDirPath, err := config.GetConfigDirPath()
	if err != nil {
		return err
	}

	repositoryDir := GetRepositoryDirPath(configDirPath, repositoryURL)

	gitDir := filepath.Join(repositoryDir, ".git")

	if _, err := os.Stat(gitDir); err == nil {
		fmt.Printf("El repositorio ya está clonado en '%s'. Omitiendo la clonación.\n", repositoryDir)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error al verificar el directorio .git: %w", err)
	}

	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio '%s': %w", configDirPath, err)
	}

	cmd := exec.Command("git", "clone", repositoryURL, repositoryDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Clonando repositorio '%s' en '%s'...\n", repositoryURL, repositoryDir)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al ejecutar git clone: %w", err)
	}

	return nil
}

func GetRepositoryDirPath(configDirPath string, repositoryURL string) string {
	repositoryName := extractRepositoryName(repositoryURL)
	return filepath.Join(configDirPath, repositoryName)
}

func extractRepositoryName(repositoryURL string) string {
	parts := strings.Split(repositoryURL, "/")
	fullName := parts[len(parts)-1]
	return strings.TrimSuffix(fullName, ".git")
}
