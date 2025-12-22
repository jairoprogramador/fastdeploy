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

func NewWorkspace(
	rootPath vos.RootPath,
	projectName vos.ProjectName,
	templateName vos.TemplateName) (*Workspace, error) {

	return &Workspace{
		rootPath:     rootPath,
		projectName:  projectName,
		templateName: templateName,
	}, nil
}

func (w *Workspace) TemplatePath() string {
	return filepath.Join(w.rootPath.Path(), "repositories", w.templateName.String())
}

func (w *Workspace) StepTemplatePath(stepName string) string {
	return filepath.Join(w.TemplatePath(), "steps", stepName)
}

func (w *Workspace) VarsTemplatePath(stepName, environment string) string {
	return filepath.Join(w.TemplatePath(), "variables", environment, stepName)
}

func (w *Workspace) WorkspacePath() string {
	return filepath.Join(w.rootPath.Path(), w.projectName.String(), w.templateName.String())
}

func (w *Workspace) VarsDirPath() string {
	return filepath.Join(w.WorkspacePath(), "vars")
}

func (w *Workspace) VarsFilePath(fileName vos.FileName) string {
	return filepath.Join(w.VarsDirPath(), fileName.String())
}

func (w *Workspace) WorkdirPath() string {
	return filepath.Join(w.WorkspacePath(), "workdir")
}

func (w *Workspace) ScopeWorkdirPath(scope string, stepName string) string {
	return filepath.Join(w.WorkdirPath(), scope, stepName)
}

func (w *Workspace) StateDirPath() string {
	return filepath.Join(w.WorkspacePath(), "state")
}

func (w *Workspace) StateTablePath(stateName string) (string, error) {
	fileName, err := vos.NewStateFileName(stateName)
	if err != nil {
		return "", err
	}
	return filepath.Join(w.StateDirPath(), fileName.String()), nil
}
