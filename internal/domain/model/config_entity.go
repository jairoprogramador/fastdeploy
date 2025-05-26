package model

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/constant"
)

type ConfigEntity struct {
	Organization string `yaml:"organization"`
	TeamName     string `yaml:"teamName"`
}

func NewConfigEntity() *ConfigEntity {
	return &ConfigEntity{
		Organization: constant.DefaultOrganization,
		TeamName:     constant.DefaultTeamName,
	}
}

func (p *ConfigEntity) IsComplete() bool {
	return p.Organization != "" && p.TeamName != ""
}
