package aggregates

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/vos"

type Config struct {
	project    vos.Project
	template   vos.Template
	technology vos.Technology
	runtime    vos.Runtime
	state      vos.State
}

func NewConfig(
	project vos.Project,
	template vos.Template,
	technology vos.Technology,
	runtime vos.Runtime,
	state vos.State) *Config {

	return &Config{
		project:    project,
		template:   template,
		technology: technology,
		runtime:    runtime,
		state:      state,
	}
}

func (c Config) Project() vos.Project { return c.project }
func (c Config) Template() vos.Template { return c.template }
func (c Config) Technology() vos.Technology { return c.technology }
func (c Config) Runtime() vos.Runtime { return c.runtime }
func (c Config) State() vos.State { return c.state }