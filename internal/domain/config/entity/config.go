package entity

import (
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
)

// ConfigEntity represents the global configuration of the application
type ConfigEntity struct {
	Organization string `yaml:"organization"`
	TeamName     string `yaml:"teamName"`
}

// NewConfigEntity creates a new ConfigEntity with default values
func NewConfigEntity() *ConfigEntity {
	return &ConfigEntity{
		Organization: constant.DefaultOrganization,
		TeamName:     constant.DefaultTeamName,
	}
}

// IsComplete checks if all required fields are set
func (p *ConfigEntity) IsComplete() bool {
	return p.Organization != "" && p.TeamName != ""
}
