package project

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/repository"
)

const (
	defaultVersion       = "1.0.0"
	defaultTeamName      = "itachi"
	defaultOrganization  = "FastDeploy"
	defaultTechnology    = "springboot"
	defaultRepositoryURL = "https://github.com/jairoprogramador/mydeploy.git"
)

type Initializer struct{}

func NewInitializer() *Initializer {
	return &Initializer{}
}

func (i *Initializer) InitializeProject() (*ProjectEntity, error) {
	projectName, err := getProjectName()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	projectID, err := generateUniqueID(projectName)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID único: %w", err)
	}

	configEntity, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("error al cargar la configuración global: %w", err)
	}

	org := defaultOrganization
	repo := defaultRepositoryURL
	team := defaultTeamName

	if configEntity.Organization != "" {
		org = configEntity.Organization
	}
	if configEntity.Repository != "" {
		repo = configEntity.Repository
	}
	if configEntity.TeamName != "" {
		team = configEntity.TeamName
	}

	projectEntity := &ProjectEntity{
		Organization: org,
		ProjectID:    projectID,
		ProjectName:  projectName,
		Repository:   repo,
		Technology:   defaultTechnology,
		Version:      defaultVersion,
		TeamName:     team,
	}

	if err := repository.CloneRepository(repo); err != nil {
		return nil, fmt.Errorf("error al clonar el repositorio: %w", err)
	}

	if err := Save(*projectEntity); err != nil {
		return nil, err
	}

	return projectEntity, nil
}

func CheckIfAlreadyInitialized() bool {
	info, err := os.Stat(constants.ProjectFileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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
