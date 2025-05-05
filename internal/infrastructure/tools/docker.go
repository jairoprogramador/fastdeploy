package tools

import (
	constants "deploy/internal/domain"
	"strings"
	"text/template"
)

func GetImageId(hashCommit string) (string, error) {
	imageId, err := ExecuteCommand("docker", "images", "-q", hashCommit)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(imageId), nil
}

func GetContainersId(imageId string) ([]string, error) {
	ancestor := "ancestor=" + imageId
	containerIds, err := ExecuteCommand("docker", "ps", "-qa", "--filter", ancestor)
	if err != nil {
		return []string{}, err
	}

	if containerIds == "" {
		return []string{}, nil
	}

	return getArray(containerIds), nil
}

func SonarScanner(token, projectKey, projectName, projectPath, cacheDir, tmpDir, scannerWorkDir, sourcePath, testPath, binaryPath, testBinaryPath string) error {
	args := []string{
		"run",
		"--rm",
		"--network=host",
		"-v", tmpDir + ":/opt/sonar-scanner/.sonar/_tmp",
		"-e", "SONAR_HOST_URL=http://localhost:9000",
		"-e", "SONAR_SCANNER_OPTS=-Xmx1024m -Djava.io.tmpdir=/opt/sonar-scanner/.sonar/_tmp",
		"-v", projectPath + ":/usr/src",
		"-v", cacheDir + ":/opt/sonar-scanner/.sonar/cache",
		"-v", scannerWorkDir + ":/opt/sonar-scanner/.scannerwork",
		"sonarsource/sonar-scanner-cli:latest",
		"-Dsonar.token=" + token,
		"-Dsonar.projectKey=" + projectKey,
		"-Dsonar.projectName=" + projectName,
		"-Dsonar.sources=" + sourcePath,
		"-Dsonar.tests=" + testPath,
		"-Dsonar.java.binaries=" + binaryPath,          //solo para java
		"-Dsonar.java.test.binaries=" + testBinaryPath, //solo para java
		"-Dsonar.sourceEncoding=UTF-8",
		"-Dsonar.scm.provider=git",
		"-Dsonar.tempFolder=/opt/sonar-scanner/.sonar/_tmp",
		"-Dsonar.working.directory=/opt/sonar-scanner/.scannerwork",
	}
	_, err := ExecuteCommand("docker", args...)
	return err
}

func GetPortContainer(containerId string) (string, error) {
	return ExecuteCommand("docker", "port", containerId)
}

func BuildImage(hashCommit string, filePath string) error {
	_, err := ExecuteCommand("docker", "build", "-t", hashCommit, "-f", filePath, ".")
	return err
}

func BuildContainer(filePath string) error {
	_, err := ExecuteCommand("docker", "compose", "-f", filePath, "up", "-d")
	return err
}

func Start(containerId string) error {
	_, err := ExecuteCommand("docker", "start", containerId)
	return err
}

func Restart(containerId string) error {
	_, err := ExecuteCommand("docker", "restart", containerId)
	return err
}

func GetContainerStatus(containerId string) (string, error) {
	status, err := ExecuteCommand("docker", "inspect", "--format", "{{.State.Status}}", containerId)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(status), nil
}

func StartContainerIfStopped(containerId string) error {
	status, err := GetContainerStatus(containerId)
	if err != nil {
		return err
	}
	if status != "running" {
		if err := Start(containerId); err != nil {
			return err
		}
	}
	return nil
}

func GetDockerfileContent(param map[string]string, filePath string) (string, error) {
	type DockerParams struct {
		FileName      string
		CommitMessage string
		CommitHash    string
		CommitAuthor  string
		Team          string
		Organization  string
	}

	params := DockerParams{
		FileName:      param[constants.FileNameKey],
		CommitMessage: param[constants.CommitMessageKey],
		CommitHash:    param[constants.CommitHashKey],
		CommitAuthor:  param[constants.CommitAuthorKey],
		Team:          param[constants.TeamKey],
		Organization:  param[constants.OrganizationKey],
	}

	dockerfileTemplate, err := template.ParseFiles(filePath)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = dockerfileTemplate.Execute(&result, params)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func GetComposeContent(param map[string]string, filePath string) (string, error) {

	type DockerParams struct {
		NameDelivery string
		CommitHash   string
		Port         string
	}

	params := DockerParams{
		NameDelivery: param[constants.NameDeliveryKey],
		CommitHash:   param[constants.CommitHashKey],
		Port:         param[constants.PortKey],
	}

	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = tmpl.Execute(&result, params)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func GetSonarqubeComposeContent(homeDir, templateData string) (string, error) {

	type ComposeParams struct {
		HomeDir string
	}

	params := ComposeParams{
		HomeDir: homeDir,
	}

	//tmpl, err := template.ParseFiles(filePath)
	tmpl, err := template.New("compose").Parse(templateData)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = tmpl.Execute(&result, params)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func getArray(data string) []string {
	return strings.Split(strings.TrimSpace(data), "\n")
}
