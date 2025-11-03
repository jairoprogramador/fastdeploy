package aggregates

import (
	proEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/entities"
	proVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/vos"
)

type MyProject struct {
	project  *proEnt.Project
	template proVos.Template
	state    proVos.State
}

func NewMyProject(
	project *proEnt.Project,
	template proVos.Template,
	state proVos.State) *MyProject {

	return &MyProject{
		project:  project,
		template: template,
		state:    state,
	}
}

func (c *MyProject) Project() *proEnt.Project  { return c.project }
func (c *MyProject) Template() proVos.Template { return c.template }
func (c *MyProject) State() proVos.State       { return c.state }

func (c *MyProject) WithProject(newProject *proEnt.Project) *MyProject {
	return &MyProject{
		project:  newProject,
		template: c.template,
		state:    c.state,
	}
}
