package project

import (
	"crypto/sha1"
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultVersion       = "1.0.0"
	defaultTeamName      = "itachi"
	defaultYAMLFile      = "deploy.yaml"
	defaultOrganization  = "FastDeploy"
	defaultRepositoryURL = "https://github.com/jairoprogramador/mydeploy.git"
)

type Initializer struct{}

func NewInitializer() *Initializer {
	return &Initializer{}
}

func CheckIfAlreadyInitialized() bool {
	info, err := os.Stat(defaultYAMLFile)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (i *Initializer) InitializeProject(projectName string) (*config.Config, error) {
	projectID, err := generateUniqueID(projectName)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID único: %w", err)
	}

	cfg := &config.Config{
		Organization: defaultOrganization,
		ProjectID:    projectID,
		ProjectName:  projectName,
		Repository:   defaultRepositoryURL,
		Technology:   "java",
		Version:      defaultVersion,
		TeamName:     defaultTeamName,
	}

	if err := i.cloneRepository(cfg.Repository); err != nil {
		return nil, fmt.Errorf("error al clonar el repositorio: %w", err)
	}

	yamlData, err := cfg.ToYAML()
	if err != nil {
		return nil, fmt.Errorf("error al serializar a YAML: %w", err)
	}

	if err := os.WriteFile(defaultYAMLFile, yamlData, 0644); err != nil {
		return nil, fmt.Errorf("error al escribir el archivo YAML: %w", err)
	}

	return cfg, nil
}

func generateUniqueID(projectName string) (string, error) {
	timestamp := time.Now().String()

	data := []byte(projectName + timestamp)

	hash := sha1.New()

	_, err := hash.Write(data)
	if err != nil {
		return "", fmt.Errorf("error al generar el hash SHA-1: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (i *Initializer) cloneRepository(repoURL string) error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("no se pudo obtener el directorio de inicio del usuario: %w", err)
	}

	dotDir := filepath.Join(currentUser.HomeDir, ".fastdeploy")

	repoName := extractRepoName(repoURL)
	projectRepoDir := filepath.Join(dotDir, repoName)
	gitDir := filepath.Join(projectRepoDir, ".git")

	if _, err := os.Stat(gitDir); err == nil {
		fmt.Printf("El repositorio ya está clonado en '%s'. Omitiendo la clonación.\n", projectRepoDir)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error al verificar el directorio .git: %w", err)
	}

	if err := os.MkdirAll(dotDir, 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio '%s': %w", dotDir, err)
	}

	cmd := exec.Command("git", "clone", repoURL, projectRepoDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Clonando repositorio '%s' en '%s'...\n", repoURL, projectRepoDir)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al ejecutar git clone: %w", err)
	}

	return nil
}

func extractRepoName(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	fullName := parts[len(parts)-1]
	return strings.TrimSuffix(fullName, ".git")
}
