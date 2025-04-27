package repository

import (
	"sync"
	"net"
	"fmt"
	"strings"
	"unicode/utf8"
	"deploy/internal/infrastructure/filesystem"
	"deploy/internal/infrastructure/tools"
	"deploy/internal/interface/presenter"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/model"
	"deploy/internal/domain"
)

type publishRepositoryImpl struct {}

var (
    instancePublishRepository     repository.PublishRepository
    instanceOncePublishRepository sync.Once
)

func GetPublishRepository() repository.PublishRepository {
    instanceOncePublishRepository.Do(func() {
        instancePublishRepository = &publishRepositoryImpl{}
    })
    return instancePublishRepository
}

func (s *publishRepositoryImpl) Prepare() *model.Response {
	_, err := tools.GetMavenVersion()
	if err != nil {
		return model.GetNewResponseError(err)
	}
	thereAreChanges, err := tools.ThereAreChanges()
	if err != nil {
		return model.GetNewResponseError(err)
	}
	if thereAreChanges {
		return model.GetNewResponseError(fmt.Errorf(constants.MessageThereAreUnconfirmedChanges))
	}
	return model.GetNewResponseMessage("")
}

func (s *publishRepositoryImpl) Build() *model.Response {
	commitHash, err := tools.GetCommitHash()
	if err != nil {
		return model.GetNewResponseError(err)
	}

	imageId, err := tools.GetImageId(commitHash)
	if err != nil {
		return model.GetNewResponseError(err)
	}

	if !imageExists(imageId) {
		err := buildProyect()
		if err != nil {
			return model.GetNewResponseError(err)
		}
	}

	response := model.GetNewResponse()
	response.SetCommitHash(commitHash)
	response.SetImageId(imageId)

	presenter.ShowSuccess("Build")
	return response
}

func (s *publishRepositoryImpl) Package(response *model.Response) *model.Response {
	if !imageExists(response.GetImageId()) {
		if err := packageInImage(response); err != nil {
			return model.GetNewResponseError(err)
		}

		imageId, err := tools.GetImageId(response.GetCommitHash())
		if err != nil {
			return model.GetNewResponseError(err)
		}
		response.SetImageId(imageId)
	}
	presenter.ShowSuccess("Package")
	return response
}

func (s *publishRepositoryImpl) Deliver(response *model.Response) *model.Response {
	containerIds, err := tools.GetContainersId(response.GetImageId())
	if err != nil {
		return model.GetNewResponseError(err)
	}

	if containerExists(containerIds) {
		for _, containerId := range containerIds {
			if err := tools.Restart(containerId); err != nil {
				return model.GetNewResponseError(err)
			}
		}
	} else {
		if err = buildContainer(response); err != nil {
			return model.GetNewResponseError(err)
		}
	}
	presenter.ShowSuccess("Deliver")
	return response
}

func (s *publishRepositoryImpl) Validate(response *model.Response) *model.Response {
	containerIds, err := tools.GetContainersId(response.GetImageId())
	if err != nil {
		return model.GetNewResponseError(err)
	}

	if containerExists(containerIds) {
		urlsContainers, err := getUrlsContainer(containerIds)
		if err != nil {
			return model.GetNewResponseError(err)
		}
		return model.GetNewResponseMessage(urlsContainers)
	}else{
		return model.GetNewResponseError(fmt.Errorf(constants.MessageErrorCreatingContainer))
	}
}

func getUrlsContainer(containerIDs []string) (string, error) {
	var result strings.Builder
	for index, id := range containerIDs {
		port, err := getHostPort(id)
		if err != nil {
			return "", err
		}
		url := fmt.Sprintf(constants.MessageSuccessPublish, index, port)
		result.WriteString(url + "\n")
	}
	return result.String(), nil	
}

func imageExists(imageId string) bool {
	return utf8.RuneCountInString(imageId) > 0
}

func containerExists(containerIds []string) bool{
	return len(containerIds) > 0
}

func buildProyect() error {
	_, err := tools.CleanAndPackage()
	if err != nil {
		return err
	}
	_, err = tools.SearchFile("target")
	if err != nil {
		return  err
	}
	return nil
}

func packageInImage(response *model.Response) error {
	commitHash := response.Data[constants.CommitHashKey]

	commitMessage, err := tools.GetCommitMessage(commitHash)
	if err != nil {
		return err
	}

	commitAuthor, err := tools.GetCommitAuthor(commitHash)
	if err != nil {
		return err
	}

	archivosJar, err := tools.SearchFile("target")
	if err != nil {
		return err
	}

	projectRepository := GetProjectRepository()
	project, err := projectRepository.Load()
	if err != nil {
		return err
	}
	
	param := make(map[string]string, 6)
	param[constants.FileNameKey] = archivosJar[0]
	param[constants.CommitHashKey] = commitHash
	param[constants.CommitMessageKey] = commitMessage
	param[constants.CommitAuthorKey] = commitAuthor
	param[constants.TeamKey] = project.TeamName
	param[constants.OrganizationKey] = project.Organization

	dockerfileContent, err := tools.GetDockerfileContent(param, constants.DockerfileTemplateFilePath)
	if err != nil {
		return err
	}
	err = filesystem.WriteFile(constants.DockerfileFilePath, dockerfileContent)
	if err != nil {
		return err
	}
	err = tools.BuildImage(commitHash, constants.DockerfileFilePath)
	if err != nil {
		return err
	}
	return filesystem.Removefile(constants.DockerfileFilePath)
}

func buildContainer(response *model.Response) error {
	commitHash := response.Data[constants.CommitHashKey]

	projectRepository := GetProjectRepository()
	project, err := projectRepository.Load()
	if err != nil {
		return err
	}
	
	param := make(map[string]string, 3)
	param[constants.NameDeliveryKey] = project.ProjectId + commitHash
	param[constants.CommitHashKey] = commitHash
	param[constants.PortKey] = getPortForContainer()

	composeContent, err := tools.GetComposeContent(param, constants.DockercomposeTemplateFilePath)
	if err != nil {
		return err
	}
	err = filesystem.WriteFile(constants.DockercomposeFilePath, composeContent)
	if err != nil {
		return err
	}
	err = tools.BuildContainer(constants.DockercomposeFilePath)
	if err != nil {
		return err
	}
	
	return filesystem.Removefile(constants.DockercomposeFilePath)
}

func getPortForContainer() string {
	startPort := 2000
	endPort := 3000

	portFree := 2000

	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", address)
		if err == nil {
			portFree = port
			ln.Close()
			break
		}
	}
	return fmt.Sprintf("%d", portFree)
}

func getHostPort(containerID string) (string, error) {
	ports, err := tools.GetPortContainer(containerID)
	if err != nil {
		return "", err
	}

	lines := strings.Split(ports, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "->")
		if len(parts) == 2 {
			hostPart := strings.TrimSpace(parts[1])
			if strings.Contains(hostPart, ":") {
				hostPort := strings.Split(hostPart, ":")[1]
				return hostPort, nil
			}
		}
	}
	return "", fmt.Errorf(constants.MessageErrorNoPortHost)
}