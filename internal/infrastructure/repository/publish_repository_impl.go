package repository

import (
	constants "deploy/internal/domain"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/infrastructure/filesystem"
	"deploy/internal/infrastructure/tools"
	"deploy/internal/interface/presenter"
	"fmt"
	"net"
	"path/filepath"
	"sync"
	"time"
	"unicode/utf8"
)

type publishRepositoryImpl struct{
	dockerRepo repository.DockerRepository
}

var (
	instancePublishRepository     repository.PublishRepository
	instanceOncePublishRepository sync.Once
)

func GetPublishRepository() repository.PublishRepository {
	instanceOncePublishRepository.Do(func() {
		instancePublishRepository = &publishRepositoryImpl{
			dockerRepo: tools.NewDockerService(),
		}
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

	imageId, err := s.dockerRepo.GetImageID(commitHash)
	if err != nil {
		return model.GetNewResponseError(err)
	}

	if !imageExists(imageId) {
		err := s.buildProyect()
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
		if err := s.packageInImage(response); err != nil {
			return model.GetNewResponseError(err)
		}
		imageId, err := s.dockerRepo.GetImageID(response.GetCommitHash())
		if err != nil {
			return model.GetNewResponseError(err)
		}
		response.SetImageId(imageId)
	}
	presenter.ShowSuccess("Package")
	return response
}

func (s *publishRepositoryImpl) Deliver(response *model.Response) *model.Response {
	containerIds, err := s.dockerRepo.GetContainersID(response.GetImageId())
	if err != nil {
		return model.GetNewResponseError(err)
	}

	if containerExists(containerIds) {
		for _, containerId := range containerIds {
			if err := s.dockerRepo.StartContainerIfStopped(containerId); err != nil {
				return model.GetNewResponseError(err)
			}
		}
	} else {
		if err = s.buildContainer(response); err != nil {
			return model.GetNewResponseError(err)
		}
	}
	presenter.ShowSuccess("Deliver")
	return response
}

func (s *publishRepositoryImpl) Validate(response *model.Response) *model.Response {
	containerIds, err := s.dockerRepo.GetContainersID(response.GetImageId())
	if err != nil {
		return model.GetNewResponseError(err)
	}

	if containerExists(containerIds) {
		urlsContainers, err := s.dockerRepo.GetUrlsContainer(containerIds)
		if err != nil {
			return model.GetNewResponseError(err)
		}
		return model.GetNewResponseMessage(urlsContainers)
	} else {
		return model.GetNewResponseError(fmt.Errorf(constants.MessageErrorCreatingContainer))
	}
}

func imageExists(imageId string) bool {
	return utf8.RuneCountInString(imageId) > 0
}

func containerExists(containerIds []string) bool {
	return len(containerIds) > 0
}

func (s *publishRepositoryImpl) buildProyect() error {
	_, err := tools.CleanAndPackage()
	if err != nil {
		return err
	}

	err = s.SonarQube()
	if err != nil {
		return err
	}

	_, err = tools.SearchFile("target")
	if err != nil {
		return err
	}
	return nil
}

func (s *publishRepositoryImpl) SonarQube() error {
	sonarqubeRepository := GetSonarqubeRepository()
	error := sonarqubeRepository.RevokeToken()
	if error != nil {
		return error
	}

	projectRepository := NewProjectRepository()
	project, err := projectRepository.Load()
	if err != nil {
		return err
	}

	token, err := sonarqubeRepository.CreateToken(project.ProjectID)
	if err != nil {
		return err
	}

	homeDir, err := filesystem.GetHomeDirectory()
	if err != nil {
		return err
	}

	cacheScannerDir := filepath.Join(homeDir, ".fastdeploy", "scanner", "cache")
	err = filesystem.RecreateDirectory(cacheScannerDir)
	if err != nil {
		return err
	}

	tmpScannerDir := filepath.Join(homeDir, ".fastdeploy", "scanner", "tmp")
	err = filesystem.RecreateDirectory(tmpScannerDir)
	if err != nil {
		return err
	}

	scannerWorkDir := filepath.Join(homeDir, ".fastdeploy", "scanner", "work")
	err = filesystem.RecreateDirectory(scannerWorkDir)
	if err != nil {
		return err
	}

	projectDirectory, err := filesystem.GetProjectDirectory()
	if err != nil {
		return err
	}

	projectKey := project.ProjectID
	projectName := project.ProjectID
	projectPath := projectDirectory
	sourcePath := "src/main"
	testPath := "src/test"
	binaryPath := "target/classes"
	testBinaryPath := "target/test-classes"

	err = s.dockerRepo.SonarScanner(token, projectKey, projectName, projectPath, cacheScannerDir, tmpScannerDir, scannerWorkDir, sourcePath, testPath, binaryPath, testBinaryPath)
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	status, err := sonarqubeRepository.GetQualityGateStatus(projectKey)
	if err != nil {
		return err
	}

	if status != "OK" && status != "PASSED" {
		return fmt.Errorf("los Quality Gates no han superado los mínimos aceptables. Estado: %s", status)
	}

	return nil
}

func (s *publishRepositoryImpl) packageInImage(response *model.Response) error {
	commitHash := response.GetCommitHash()
	fmt.Println("c dd " + commitHash)
	commitMessage, err := tools.GetCommitMessage(commitHash)
	if err != nil {
		fmt.Println("d")
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

	projectRepository := NewProjectRepository()
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

	dockerfileContent, err := s.dockerRepo.GetDockerfileContent(param, constants.DockerfileTemplateFilePath)
	if err != nil {
		return err
	}
	err = filesystem.WriteFile(constants.DockerfileFilePath, dockerfileContent)
	if err != nil {
		return err
	}
	err = s.dockerRepo.BuildImage(commitHash, constants.DockerfileFilePath)
	if err != nil {
		return err
	}
	return filesystem.RemoveFile(constants.DockerfileFilePath)
}

func (s *publishRepositoryImpl) buildContainer(response *model.Response) error {
	commitHash := response.Data[constants.CommitHashKey]

	projectRepository := NewProjectRepository()
	project, err := projectRepository.Load()
	if err != nil {
		return err
	}

	param := make(map[string]string, 3)
	param[constants.NameDeliveryKey] = project.ProjectID + commitHash[:5]
	param[constants.CommitHashKey] = commitHash
	param[constants.PortKey] = getPortForContainer()

	composeContent, err := s.dockerRepo.GetComposeContent(param, constants.DockercomposeTemplateFilePath)
	if err != nil {
		return err
	}
	err = filesystem.WriteFile(constants.DockercomposeFilePath, composeContent)
	if err != nil {
		return err
	}
	err = s.dockerRepo.BuildContainer(constants.DockercomposeFilePath)
	if err != nil {
		return err
	}

	return filesystem.RemoveFile(constants.DockercomposeFilePath)
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
