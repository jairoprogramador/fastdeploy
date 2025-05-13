package repository

import (
	"bytes"
	"context"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/template"
	"deploy/internal/infrastructure/filesystem"
	"deploy/internal/infrastructure/tools"
	"deploy/internal/interface/presenter"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
	"time"
)

const (
	sonarURL          = "http://localhost:9000"
	user              = "admin"
	passwordDefault   = "admin"
	password          = "F@stDeploy2025"
	tokenName         = "fastdeploy"
	statusAPI         = sonarURL + "/api/system/status"
	tokenGenerateAPI  = sonarURL + "/api/user_tokens/generate"
	tokenRevokeAPI    = sonarURL + "/api/user_tokens/revoke"
	changePasswordAPI = sonarURL + "/api/users/change_password"
	createProjectAPI  = sonarURL + "/api/projects/create"
	searchProjectAPI  = sonarURL + "/api/projects/search"
	projectStatusAPI  = sonarURL + "/api/qualitygates/project_status"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type sonarqubeRepositoryImpl struct {
	client            *http.Client
	user              string
	passwordDefault   string
	password          string
	tokenName         string
	changePasswordAPI string
	tokenGenerateAPI  string
	tokenRevokeAPI    string
	createProjectAPI  string
	searchProjectAPI  string
	projectStatusAPI  string
	dockerRepo        repository.DockerRepository
}

var (
	instanceSonarqubeRepository     repository.SonarqubeRepository
	instanceOnceSonarqubeRepository sync.Once
)

func GetSonarqubeRepository() repository.SonarqubeRepository {
	instanceOnceSonarqubeRepository.Do(func() {
		instanceSonarqubeRepository = &sonarqubeRepositoryImpl{
			client: &http.Client{
				Timeout: 10 * time.Second,
			},
			user:              user,
			passwordDefault:   passwordDefault,
			tokenName:         tokenName,
			password:          password,
			changePasswordAPI: changePasswordAPI,
			tokenGenerateAPI:  tokenGenerateAPI,
			tokenRevokeAPI:    tokenRevokeAPI,
			createProjectAPI:  createProjectAPI,
			searchProjectAPI:  searchProjectAPI,
			projectStatusAPI:  projectStatusAPI,
			dockerRepo:        tools.NewDockerService(),
		}
	})
	return instanceSonarqubeRepository
}

func (s *sonarqubeRepositoryImpl) Add() *model.Response {
	homeDir, err := filesystem.GetHomeDirectory()
	if err != nil {
		return model.GetNewResponseError(err)
	}

	dataDir := filepath.Join(homeDir, ".fastdeploy", "sonarqube", "volumes", "data")
	logsDir := filepath.Join(homeDir, ".fastdeploy", "sonarqube", "volumes", "logs")
	extensionsDir := filepath.Join(homeDir, ".fastdeploy", "sonarqube", "volumes", "extensions")

	err = filesystem.CreateDirectory(dataDir)
	if err != nil {
		return model.GetNewResponseError(err)
	}

	err = filesystem.CreateDirectory(logsDir)
	if err != nil {
		return model.GetNewResponseError(err)
	}

	err = filesystem.CreateDirectory(extensionsDir)
	if err != nil {
		return model.GetNewResponseError(err)
	}

	sonarqubeComposePath := filepath.Join(homeDir, ".fastdeploy", "sonarqube", "compose.yaml")
	existsDockerfile, _ := filesystem.ExistsFile(sonarqubeComposePath)

	if !existsDockerfile {
		composeContent, err := s.dockerRepo.GetSonarqubeComposeContent(homeDir, template.ComposeSonarqubeTemplate)
		if err != nil {
			return model.GetNewResponseError(err)
		}
		err = filesystem.WriteFile(sonarqubeComposePath, composeContent)
		if err != nil {
			return model.GetNewResponseError(err)
		}

		err = s.dockerRepo.BuildContainer(sonarqubeComposePath)
		if err != nil {
			return model.GetNewResponseError(err)
		}
	} else {
		containerIds, err := s.dockerRepo.GetContainersID("sonarqube:community")
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
			err = s.dockerRepo.BuildContainer(sonarqubeComposePath)
			if err != nil {
				return model.GetNewResponseError(err)
			}
		}
	}
	return s.validate()
}

func (s *sonarqubeRepositoryImpl) validate() *model.Response {
	containerIds, err := s.dockerRepo.GetContainersID("sonarqube:community")
	if err != nil {
		return model.GetNewResponseError(err)
	}

	if containerExists(containerIds) {
		urlsContainers, err :=  s.dockerRepo.GetUrlsContainer(containerIds)
		if err != nil {
			return model.GetNewResponseError(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		err = s.WaitSonarqube(ctx, 15, 50*time.Second)
		if err != nil {
			return model.GetNewResponseError(err)
		}

		projectRepository := GetProjectRepository()
		project, err := projectRepository.Load()
		if err != nil {
			return model.GetNewResponseError(err)
		}

		exists, err := s.ProjectExists(project.ProjectID)
		if err != nil {
			return model.GetNewResponseError(err)
		}
		if !exists {
			err = s.CreateProject(project.ProjectID, project.Name)
			if err != nil {
				return model.GetNewResponseError(err)
			}
		}

		if !s.CanLogin(s.user, s.password) {
			_, err := s.ChangePassword()
			if err != nil {
				return model.GetNewResponseError(err)
			}
		}

		presenter.ShowSuccess("Add sonarqube")

		urlsContainers = urlsContainers + " user=" + s.user + " password=" + s.password
		return model.GetNewResponseMessage(urlsContainers)
	} else {
		return model.GetNewResponseError(fmt.Errorf(constant.MessageErrorCreatingSonarqube))
	}
}

func (s *sonarqubeRepositoryImpl) WaitSonarqube(ctx context.Context, maxRetries int, interval time.Duration) error {
	fmt.Println("Esperando a que SonarQube esté listo...")

	for attempt := 1; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		resp, err := http.Get(statusAPI)
		if err == nil {
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Intento %d: respuesta inesperada: %d\n", attempt, resp.StatusCode)
			} else {
				var result struct {
					Status string `json:"status"`
				}

				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					fmt.Printf("Intento %d: error parseando JSON: %v\n", attempt, err)
				} else if result.Status == "UP" {
					fmt.Println("SonarQube está listo ✅")
					return nil
				}
			}
		}

		fmt.Printf("SonarQube no está listo aún (intento %d/%d). Esperando %s...\n", attempt, maxRetries, interval)
		time.Sleep(interval)
	}

	return fmt.Errorf("sonarqube no estuvo listo después de %d intentos", maxRetries)
}

func (s *sonarqubeRepositoryImpl) CreateToken(projectKey string) (string, error) {
	data := url.Values{}
	data.Set("name", s.tokenName)
	data.Set("type", "PROJECT_ANALYSIS_TOKEN")
	data.Set("projectKey", projectKey)

	body, status, err := s.sendRequest("POST", s.tokenGenerateAPI, data, s.user, s.password)
	if err != nil {
		return "", err
	}

	if status != http.StatusOK {
		return "", fmt.Errorf("error creando token: status %d, body: %s", status, string(body))
	}

	var result TokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parseando respuesta: %w", err)
	}

	return result.Token, nil
}

func (s *sonarqubeRepositoryImpl) RevokeToken() error {
	data := url.Values{}
	data.Set("name", s.tokenName)

	body, status, err := s.sendRequest("POST", s.tokenRevokeAPI, data, s.user, s.password)
	if err != nil {
		return err
	}

	if status != http.StatusNoContent && status != http.StatusNotFound {
		return fmt.Errorf("error eliminando token: status %d, body: %s", status, string(body))
	}

	return nil
}

func (s *sonarqubeRepositoryImpl) ChangePassword() (string, error) {
	data := url.Values{}
	data.Set("login", s.user)
	data.Set("previousPassword", s.passwordDefault)
	data.Set("password", s.password)

	body, status, err := s.sendRequest("POST", s.changePasswordAPI, data, s.user, s.passwordDefault)
	if err != nil {
		return "", err
	}

	if status != http.StatusOK && status != http.StatusNoContent {
		return "", fmt.Errorf("error cambiando contraseña: status %d, body: %s", status, string(body))
	}

	return s.password, nil
}

func (s *sonarqubeRepositoryImpl) CreateProject(projectKey, projectName string) error {
	data := url.Values{}
	data.Set("name", projectName)
	data.Set("project", projectKey)
	data.Set("mainBranch", "main")
	//data.Set("newCodeDefinitionType", "PREVIOUS_VERSION")

	body, status, err := s.sendRequest("POST", s.createProjectAPI, data, s.user, s.password)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return fmt.Errorf("error creando proyecto: status %d, body: %s", status, string(body))
	}

	return nil
}

func (s *sonarqubeRepositoryImpl) CanLogin(user, pass string) bool {
	_, status, err := s.sendRequest("GET", statusAPI, nil, s.user, s.password)
	if err != nil {
		return false
	}
	return status == http.StatusOK
}

func (s *sonarqubeRepositoryImpl) ProjectExists(projectKey string) (bool, error) {
	apiURL := s.searchProjectAPI + "?projects=" + projectKey
	body, status, err := s.sendRequest("GET", apiURL, nil, s.user, s.password)
	if err != nil {
		return false, err
	}

	if status != http.StatusOK {
		return false, fmt.Errorf("error consultando proyecto: status %d", status)
	}

	var result struct {
		Components []struct {
			Key string `json:"key"`
		} `json:"components"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	for _, c := range result.Components {
		if c.Key == projectKey {
			return true, nil
		}
	}
	return false, nil
}

func (s *sonarqubeRepositoryImpl) GetQualityGateStatus(projectKey string) (string, error) {
	apiURL := s.projectStatusAPI + "?projectKey=" + projectKey

	body, status, err := s.sendRequest("GET", apiURL, nil, s.user, s.password)
	if err != nil {
		return "", err
	}

	if status != http.StatusOK {
		return "", fmt.Errorf("error consultando proyecto: status %d", status)
	}

	var result struct {
		ProjectStatus struct {
			Status string `json:"status"`
		} `json:"projectStatus"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.ProjectStatus.Status, nil
}

func (s *sonarqubeRepositoryImpl) sendRequest(method, url string, data url.Values, authUser, authPass string) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, 0, fmt.Errorf("error creando request: %w", err)
	}

	req.SetBasicAuth(authUser, authPass)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("error enviando request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	return body, resp.StatusCode, nil
}
