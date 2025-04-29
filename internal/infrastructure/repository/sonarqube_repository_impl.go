package repository

import (
	"bytes"
	"deploy/internal/domain"
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
	"path/filepath"
	"sync"
	"time"
	"net/url"
	"context"
)

const (
	sonarURL   		  = "http://localhost:9000"
	user  		  	  = "admin"
	passwordDefault   = "admin"
	password          = "F@stDeploy2025"
	tokenName  		  = "fastdeploy"
	statusAPI  		  = sonarURL + "/api/system/status"
	tokenGenerateAPI  = sonarURL + "/api/user_tokens/generate"
	tokenRevokeAPI    = sonarURL + "/api/user_tokens/revoke"
	changePasswordAPI = sonarURL + "/api/users/change_password"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type sonarqubeRepositoryImpl struct {
	client          	*http.Client
    user        	    string
    passwordDefault     string
    password         	string
    tokenName       	string
    changePasswordAPI 	string
    tokenGenerateAPI  	string
	tokenRevokeAPI      string
}

var (
    instanceSonarqubeRepository   repository.SonarqubeRepository
    instanceOnceSonarqubeRepository sync.Once
)

func GetSonarqubeRepository() repository.SonarqubeRepository {
    instanceOnceSonarqubeRepository .Do(func() {
        instanceSonarqubeRepository  = &sonarqubeRepositoryImpl{
			client: &http.Client{
				Timeout: 10 * time.Second,
			},
			user:		    	user,
			passwordDefault:    passwordDefault,
			tokenName:       	tokenName,
			password:         	password,
			changePasswordAPI: 	changePasswordAPI,
			tokenGenerateAPI:  	tokenGenerateAPI,
			tokenRevokeAPI:     tokenRevokeAPI,
		}
    })
    return instanceSonarqubeRepository 
}

func (s *sonarqubeRepositoryImpl) Add() *model.Response {
	//pass sdkfjlas33K@Dw
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
	existsDockerfile := filesystem.FileExists(sonarqubeComposePath)

	if !existsDockerfile {
		if err := filesystem.CreateDirectoryFilePath(sonarqubeComposePath); err != nil {
			return model.GetNewResponseError(err)
		}

		composeContent, err := tools.GetSonarqubeComposeContent(homeDir, template.ComposeSonarqubeTemplate)
		if err != nil {
			return model.GetNewResponseError(err)
		}
		err = filesystem.WriteFile(sonarqubeComposePath, composeContent)
		if err != nil {
			return model.GetNewResponseError(err)
		}

		err = tools.BuildContainer(sonarqubeComposePath)
		if err != nil {
			return model.GetNewResponseError(err)
		}
	} else {
		containerIds, err := tools.GetContainersId("sonarqube:community")
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
			err = tools.BuildContainer(sonarqubeComposePath)
			if err != nil {
				return model.GetNewResponseError(err)
			}
		}
	}
	return s.validate()
}

func (s *sonarqubeRepositoryImpl) validate() *model.Response {
	containerIds, err := tools.GetContainersId("sonarqube:community")
	if err != nil {
		return model.GetNewResponseError(err)
	}

	if containerExists(containerIds) {
		urlsContainers, err := getUrlsContainer(containerIds)
		if err != nil {
			return model.GetNewResponseError(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		err = s.WaitSonarqube(ctx, 15, 50*time.Second)
		if err != nil {
			return model.GetNewResponseError(err)
		}

		password, err := s.ChangePassword()
		if err != nil {
			return model.GetNewResponseError(err)
		}
		presenter.ShowSuccess("Add sonarqube")

		urlsContainers = urlsContainers + " user="+s.user+" password="+password
		return model.GetNewResponseMessage(urlsContainers)
	}else{
		return model.GetNewResponseError(fmt.Errorf(constants.MessageErrorCreatingSonarqube))
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

func (s *sonarqubeRepositoryImpl) CreateToken() (string, error) {
	data := url.Values{}
    data.Set("name", s.tokenName)

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

    if status != http.StatusOK && status != http.StatusNoContent{
        return "", fmt.Errorf("error cambiando contraseña: status %d, body: %s", status, string(body))
    }

    return s.password, nil
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
