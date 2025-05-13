package executor

import (
	irepository "deploy/internal/domain/repository"
	"deploy/internal/infrastructure/repository"
	"deploy/internal/domain/condition"
	"deploy/internal/domain/model"
	"deploy/internal/domain/variable"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/template"
	"fmt"
	"context"
)

type ContainerExecutor struct {
	BaseExecutor
	variables *variable.VariableStore
	commandRunner    CommandRunner
	conditionFactory *condition.ConditionFactory
	containerRepository irepository.ContainerRepository
	fileRepository irepository.FileRepository

}

func GetContainerExecutor(variables *variable.VariableStore) *ContainerExecutor {

	return &ContainerExecutor {
		BaseExecutor: BaseExecutor {},
		variables: variables,
		commandRunner:    GetCommandRunner(),
		conditionFactory: condition.GetConditionFactory(),
		containerRepository: repository.GetContainerRepository(),
		fileRepository: repository.GetFileRepository(),
	}
}

func (e *ContainerExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.prepareContext(ctx, step)
	defer cancel()

	return e.handleRetry(step, func() error {
		// Preparar variables locales
		e.variables.PushScope(step.Variables)
		defer e.variables.PopScope()

		// Ejecutar comando
		fmt.Printf("---------------%s-----------------\n", step.Name)

		err := e.Delete(ctx)
		if err != nil {
			return err
		}

		err = e.CreateImage(ctx)
		if err != nil {
			return err
		}

		return e.CreateContainer(ctx)
	})
}

func (e *ContainerExecutor) Delete(ctx context.Context) error {
	pathDockerCompose := e.fileRepository.GetFullPathDockerCompose(e.variables)
    if e.fileRepository.ExistsFile(pathDockerCompose) {
		command := fmt.Sprintf("docker compose -f %s down --rmi local --remove-orphans -v", pathDockerCompose)
        _, err := e.commandRunner.Run(ctx, command)
		if err != nil {
			return err
		}
    }
    return nil
}

func (e *ContainerExecutor) CreateImage(ctx context.Context) error {
	pathDockerfileTemplate := e.fileRepository.GetFullPathDockerfileTemplate(e.variables)
	if !e.fileRepository.ExistsFile(pathDockerfileTemplate) {
		err := e.containerRepository.CreateFile(pathDockerfileTemplate, template.DockerfileTemplate)
		if err != nil {
			return err
		}
	}
	pathDockerfile := e.fileRepository.GetFullPathDockerfile(e.variables)
	if e.fileRepository.ExistsFile(pathDockerfile) {
		err := e.fileRepository.DeleteFile(pathDockerfile)
		if err != nil {
			return err
		}
	}
	err := e.containerRepository.CreateDockerfile(pathDockerfile, pathDockerfileTemplate, e.variables)
	if err != nil {
		return err
	}
	return e.createImagenFromDockerfile(ctx, pathDockerfile)
}

func (e *ContainerExecutor) createImagenFromDockerfile(ctx context.Context, pathDockerfile string) error {
	commitHash := e.variables.Get(constant.VAR_COMMIT_HASH)
	projectVersion := e.variables.Get(constant.VAR_PROJECT_VERSION)

	command := fmt.Sprintf("docker build -t %s:%s -f %s .", commitHash, projectVersion, pathDockerfile)
    _, err := e.commandRunner.Run(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

func (e *ContainerExecutor) CreateContainer(ctx context.Context) error {
	pathDockerComposeTemplate := e.fileRepository.GetFullPathDockerComposeTemplate(e.variables)
	if !e.fileRepository.ExistsFile(pathDockerComposeTemplate) {
		err := e.containerRepository.CreateFile(pathDockerComposeTemplate, template.ComposeTemplate)
		if err != nil {
			return err
		}
	}		 
	pathDockerCompose := e.fileRepository.GetFullPathDockerCompose(e.variables)
	if e.fileRepository.ExistsFile(pathDockerCompose) {
		err := e.fileRepository.DeleteFile(pathDockerCompose)
		if err != nil {
			return err
		}
	}
	err := e.containerRepository.CreateDockerCompose(pathDockerCompose, pathDockerComposeTemplate, e.variables)
	if err != nil {
		return err
	}
	
	return e.createContainerFromDockerCompose(ctx, pathDockerCompose)
}

func (e *ContainerExecutor) createContainerFromDockerCompose(ctx context.Context, pathDockerCompose string) error {
	command := fmt.Sprintf("docker compose -f %s up -d", pathDockerCompose)
    _, err := e.commandRunner.Run(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

