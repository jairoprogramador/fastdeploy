package application

import (
	"deploy/internal/infrastructure/repository"
	"deploy/internal/domain/service"
	"deploy/internal/domain/executor"
	"deploy/internal/domain/variable"
)

func getProjectService() service.ProjectServiceInterface {
	yamlRepository := repository.GetYamlRepository()
	fileRepository := repository.GetFileRepository()

	globalConfigService := service.GetGlobalConfigService(yamlRepository, fileRepository)
	return service.GetProjectService(yamlRepository, globalConfigService, fileRepository)
}

func getDeploymentService() service.DeploymentServiceInterface {
	yamlRepository := repository.GetYamlRepository()
	fileRepository := repository.GetFileRepository()

	return service.GetDeploymentService(yamlRepository, fileRepository)
}

func getCommandExecutor(variableStore *variable.VariableStore) *executor.CommandExecutor {
	return executor.GetCommandExecutor(variableStore)
}

func getContainerExecutor(variableStore *variable.VariableStore) *executor.ContainerExecutor {
	containerRepository := repository.GetContainerRepository()
	fileRepository := repository.GetFileRepository()

	return executor.GetContainerExecutor(containerRepository, fileRepository, variableStore)
}

func getStoreService() service.StoreServiceInterface {
	return service.GetStoreService(projectModel)
}

