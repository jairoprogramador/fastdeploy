package service

import (
	"context"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"fmt"
	"strings"
	"sync"
)

type DockerServiceInterface interface {
	ExistsContainer(ctx context.Context, variableStore *model.VariableStore) (bool, error)
	DockerComposeUp(ctx context.Context, pathDockerCompose string, variableStore *model.VariableStore) (string, error)
	DockerComposeDownLocal(ctx context.Context, pathDockerCompose string) error
}

type DockerService struct {
	executorService ExecutorServiceInterface
}

var (
	instanceDockerService     *DockerService
	instanceOnceDockerService sync.Once
)

func GetDockerService() DockerServiceInterface {
	instanceOnceDockerService.Do(func() {
		instanceDockerService = &DockerService{
			executorService: GetExecutorService(),
		}
	})
	return instanceDockerService
}

func (d *DockerService) ExistsContainer(ctx context.Context, variableStore *model.VariableStore) (bool, error) {
	commitHash := variableStore.Get(constant.VAR_COMMIT_HASH)
	version := variableStore.Get(constant.VAR_PROJECT_VERSION)

	command := fmt.Sprintf("docker ps -aq --filter ancestor=%s:%s", commitHash, version)
	containerId, err := d.executorService.Run(ctx, command)
	if err != nil {
		return false, err
	}
	return len(containerId) > 0, nil
}

func (d *DockerService) DockerComposeUp(ctx context.Context, pathDockerCompose string, variableStore *model.VariableStore)  (string, error) {
	command := fmt.Sprintf("docker compose -f %s up -d --build", pathDockerCompose)
	_, err := d.executorService.Run(ctx, command)
	if err == nil {
		return d.getContainerURLs(ctx, variableStore)
	}
	return "", err
}

func (d *DockerService) DockerComposeDownLocal(ctx context.Context, pathDockerCompose string) error {
	command := fmt.Sprintf("docker compose -f %s down --rmi local --remove-orphans -v", pathDockerCompose)
	_, err := d.executorService.Run(ctx, command)
	return err
}

func (d *DockerService) getContainerURLs(ctx context.Context, variableStore *model.VariableStore) (string, error) {
	containerIDs, err := d.getIdsContainerUp(ctx, variableStore)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	for i, containerID := range containerIDs {
		port, err := d.getPortHostContainer(ctx, containerID)
		if err != nil {
			return "", err
		}
		url := fmt.Sprintf("service available in: http://localhost:%s/", port)
		result.WriteString(url)

		if i < len(containerIDs)-1 {
			result.WriteString("\n")
		}
	}
	return result.String(), nil
}

func (d *DockerService) getPortHostContainer(ctx context.Context, containerID string) (string, error) {
	command := fmt.Sprintf("docker port %s", containerID)
	ports, err := d.executorService.Run(ctx, command)
	if err != nil {
		return "", err
	}

	lines := strings.Split(ports, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			portHost := strings.TrimSpace(parts[1])
			if strings.TrimSpace(portHost) != "" {
				return portHost, nil
			}
		}
	}
	return "", fmt.Errorf(constant.MsgErrorNoPortHost)
}

func (d *DockerService) getIdsContainerUp(ctx context.Context, variableStore *model.VariableStore) ([]string, error) {
	commitHash := variableStore.Get(constant.VAR_COMMIT_HASH)
	version := variableStore.Get(constant.VAR_PROJECT_VERSION)

	command := fmt.Sprintf("docker ps -q --filter ancestor=%s:%s", commitHash, version)
	containerIds, err := d.executorService.Run(ctx, command)
	if err != nil {
		return []string{}, err
	}
	containerIds = strings.TrimSpace(containerIds)
	return strings.Split(containerIds, "\n"), nil
}
