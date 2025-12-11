package aggregates

import (
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/workspace/vos"
)

type Workspace struct {
	rootPath     vos.RootPath
	projectName  vos.ProjectName
	templateName vos.TemplateName
}

func NewWorkspace(rootPath vos.RootPath, projectName vos.ProjectName, templateName vos.TemplateName) (*Workspace, error) {
	return &Workspace{
		rootPath:     rootPath,
		projectName:  projectName,
		templateName: templateName,
	}, nil
}

func (w *Workspace) templatePath() string {
	return filepath.Join(w.rootPath.Path(), w.projectName.String(), w.templateName.String())
}

func (w *Workspace) VarsDirPath() string {
	return filepath.Join(w.templatePath(), "vars")
}

func (w *Workspace) VarsFilePath(fileName vos.FileName) string {
	return filepath.Join(w.VarsDirPath(), fileName.String())
}

func (w *Workspace) WorkdirPath() string {
	return filepath.Join(w.templatePath(), "workdir")
}

func (w *Workspace) ScopeWorkdirPath(scope vos.ScopeName, stepName vos.StepName) string {
	return filepath.Join(w.WorkdirPath(), scope.String(), stepName.String())
}

func (w *Workspace) StateDirPath() string {
	return filepath.Join(w.templatePath(), "state")
}

func (w *Workspace) StateTablePath(fileName vos.FileName) string {
	return filepath.Join(w.StateDirPath(), fileName.String())
}
