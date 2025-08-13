package git

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/config"
	"path/filepath"
	"strings"
)

type GitPathResolver interface {
	GetDirectoryPath(repositoryURL string) (string, error)
}

type GitPathResolverImpl struct {
	configPathResolver config.ConfigPathResolver
}

func NewGitPathResolver(configPathResolver config.ConfigPathResolver) GitPathResolver {
	return &GitPathResolverImpl{configPathResolver: configPathResolver}
}

func (gpr *GitPathResolverImpl) GetDirectoryPath(repositoryURL string) (string, error) {
	configDirPath, err := gpr.configPathResolver.GetConfigDirPath()
	if err != nil {
		return "", err
	}

	repositoryDir := gpr.getPath(configDirPath, repositoryURL)

	return repositoryDir, nil
}

func (gpr *GitPathResolverImpl) getPath(directoryPath string, repositoryURL string) string {
	repositoryName := gpr.getRepositoryName(repositoryURL)
	return filepath.Join(directoryPath, repositoryName)
}

func (gpr *GitPathResolverImpl) getRepositoryName(repositoryURL string) string {
	parts := strings.Split(repositoryURL, "/")
	fullName := parts[len(parts)-1]
	return strings.TrimSuffix(fullName, ".git")
}
