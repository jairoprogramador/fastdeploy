package tools

import "strings"
import "deploy/internal/domain"
import "text/template"

func GetImageId(hashCommit string) (string, error){
	imageId, err := ExecuteCommand("docker", "images", "-q", hashCommit)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(imageId), nil
}

func GetContainersId(imageId string) ([]string, error){
	ancestor := "ancestor=" + imageId
	containerIds, err := ExecuteCommand("docker", "ps", "-q", "--filter", ancestor)
	if err != nil {
		return []string{}, err
	}

	if containerIds == "" {
        return []string{}, nil
    }

	return getArray(containerIds), nil
}

func GetPortContainer(containerId string) (string, error){
	return ExecuteCommand("docker", "port", containerId)
}

func BuildImage(hashCommit string, filePath string) error {
	_, err := ExecuteCommand("docker", "build", "-t", hashCommit, "-f",filePath, ".")
	return err
}

func BuildContainer(filePath string) error {
	_, err := ExecuteCommand("docker", "compose", "-f", filePath, "up","-d")
	return err
}

func Restart(containerId string) error {
	_, err := ExecuteCommand("docker", "restart", containerId)
	return err
}

func GetDockerfileContent(param map[string]string, filePath string) (string, error) {
	type DockerParams struct {
		FileName string
		CommitMessage string
		CommitHash string
		CommitAuthor string
		Team string
    	Organization string
	}

	params := DockerParams{
		FileName: param[constants.FileNameKey],
		CommitMessage: param[constants.CommitMessageKey],
		CommitHash: param[constants.CommitHashKey],
		CommitAuthor: param[constants.CommitAuthorKey],
		Team: param[constants.TeamKey],
		Organization: param[constants.OrganizationKey],
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
		CommitHash string
		Port string
	}

	params := DockerParams{
		NameDelivery: param[constants.NameDeliveryKey],
		CommitHash: param[constants.CommitHashKey],
		Port: param[constants.PortKey],
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

func getArray(data string) []string{
	return strings.Split(strings.TrimSpace(data), "\n")
}