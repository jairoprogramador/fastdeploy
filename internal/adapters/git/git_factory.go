package git

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/config"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/git"
)

type GitFactory struct{}

func NewGitFactory() *GitFactory {
	return &GitFactory{}
}

func (gf *GitFactory) CreateService() domain.GitService {
	pathResolver := gf.CreatePathResolver()
	return NewGitService(pathResolver)
}

func (gf *GitFactory) CreatePathResolver() GitPathResolver {
	configPathResolver := config.NewConfigFactory().CreatePathResolver()
	return NewGitPathResolver(configPathResolver)
}
