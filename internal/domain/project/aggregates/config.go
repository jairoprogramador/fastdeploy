package aggregates

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/vos"

type Config struct {
	project  vos.Project
	template vos.Template
	state    vos.State
}

func NewConfig(
	project vos.Project,
	template vos.Template,
	state vos.State) *Config {

	return &Config{
		project:  project,
		template: template,
		state:    state,
	}
}

func (c *Config) Project() vos.Project   { return c.project }
func (c *Config) Template() vos.Template { return c.template }
func (c *Config) State() vos.State       { return c.state }
func (c *Config) SetProjectRevision(revision string) {
	c.project = c.project.WithRevision(revision)
}
