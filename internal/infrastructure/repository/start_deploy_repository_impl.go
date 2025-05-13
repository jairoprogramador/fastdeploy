package repository

import (
	"deploy/internal/domain/constant"
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

type publishRepositoryImpl struct {
	dockerRepo repository.DockerRepository
}

var (
	instancePublishRepository     repository.PublishRepository
	instanceOncePublishRepository sync.Once
)

func GetPublishRepository() repository.PublishRepository {
	instanceOncePublishRepository.Do(func() {
		instancePublishRepository = &publishRepositoryImpl {
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
		return model.GetNewResponseError(fmt.Errorf(constant.MessageThereAreUnconfirmedChanges))
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
		return model.GetNewResponseError(fmt.Errorf(constant.MessageErrorCreatingContainer))
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

	_, err = tools.GetFullPathFiles("target")
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

	projectRepository := GetProjectRepository()
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
	commitMessage, err := tools.GetCommitMessage(commitHash)
	if err != nil {
		return err
	}

	commitAuthor, err := tools.GetCommitAuthor(commitHash)
	if err != nil {
		return err
	}

	archivosJar, err := tools.GetFullPathFiles("target")
	if err != nil {
		return err
	}

	projectRepository := GetProjectRepository()
	project, err := projectRepository.Load()
	if err != nil {
		return err
	}

	param := make(map[string]string, 6)
	param[constant.FileNameKey] = archivosJar[0]
	param[constant.CommitHashKey] = commitHash
	param[constant.CommitMessageKey] = commitMessage
	param[constant.CommitAuthorKey] = commitAuthor
	param[constant.TeamKey] = project.TeamName
	param[constant.OrganizationKey] = project.Organization

	dockerfileContent, err := s.dockerRepo.GetDockerfileContent(param, constant.DockerfileTemplateFilePath)
	if err != nil {
		return err
	}
	err = filesystem.WriteFile(constant.DockerfileFilePath, dockerfileContent)
	if err != nil {
		return err
	}
	err = s.dockerRepo.BuildImage(commitHash, constant.DockerfileFilePath)
	if err != nil {
		return err
	}
	return filesystem.RemoveFile(constant.DockerfileFilePath)
}

func (s *publishRepositoryImpl) buildContainer(response *model.Response) error {
	commitHash := response.Data[constant.CommitHashKey]

	projectRepository := GetProjectRepository()
	project, err := projectRepository.Load()
	if err != nil {
		return err
	}

	param := make(map[string]string, 3)
	param[constant.NameDeliveryKey] = project.ProjectID + commitHash[:5]
	param[constant.CommitHashKey] = commitHash
	param[constant.PortKey] = getPortForContainer()

	composeContent, err := s.dockerRepo.GetComposeContent(param, constant.DockercomposeTemplateFilePath)
	if err != nil {
		return err
	}
	err = filesystem.WriteFile(constant.DockercomposeFilePath, composeContent)
	if err != nil {
		return err
	}
	err = s.dockerRepo.BuildContainer(constant.DockercomposeFilePath)
	if err != nil {
		return err
	}

	return filesystem.RemoveFile(constant.DockercomposeFilePath)
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
