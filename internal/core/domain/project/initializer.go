package project

import (
	"crypto/sha1"
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
	"os"
	"time"
)

const (
	defaultVersion      = "1.0.0"
	defaultTeamName     = "itachi"
	defaultYAMLFile     = "deploy.yaml"
	defaultOrganization = "FastDeploy"
)

type Initializer struct{}

func NewInitializer() *Initializer {
	return &Initializer{}
}

func (i *Initializer) CheckIfAlreadyInitialized() bool {
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
		Technology:   "java",
		Version:      defaultVersion,
		TeamName:     defaultTeamName,
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
