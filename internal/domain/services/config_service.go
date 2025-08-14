package services

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/configuration"
	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/project"
	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/repository"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/repositories"
)

type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// Load carga la configuración desde un archivo usando el repositorio.
// Si el repositorio indica que el archivo no existe, debe llamar a Config.create_default()
// para obtener la configuración por defecto.
func (cs *ConfigService) Load(repo repositories.ConfigRepository) (configuration.Configuration, error) {
	_, err := repo.Load()
	if err != nil {
		// Si el archivo no existe, crear configuración por defecto
		return cs.createDefault(), nil
	}

	// Aquí se convertiría el diccionario a la entidad Configuration
	// Por ahora retornamos una configuración por defecto
	return cs.createDefault(), nil
}

// Save guarda la entidad Config en un archivo usando el repositorio.
func (cs *ConfigService) Save(repo repositories.ConfigRepository, config configuration.Configuration) error {
	// Aquí se convertiría la entidad Configuration a un diccionario
	// Por ahora guardamos un diccionario vacío
	data := make(map[string]interface{})
	return repo.Save(data)
}

// createDefault crea una configuración por defecto
func (cs *ConfigService) createDefault() configuration.Configuration {
	// Crear entidades por defecto
	org, _ := project.NewOrganization("default-org")
	team, _ := project.NewTeam("default-team")
	repoURL, _ := repository.NewRepositoryURL("https://github.com/default/repo")
	repo := repository.NewRepository(repoURL)

	return configuration.NewConfiguration(org, team, repo)
}
