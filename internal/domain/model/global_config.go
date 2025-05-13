package model

import (
	"deploy/internal/domain/constant"
)

type GlobalConfig struct {
	Organization string `yaml:"organization"`
	TeamName     string `yaml:"teamName"`
}

func NewGlobalConfig(organization, teamName string) *GlobalConfig {
	if organization == "" {
		organization = constant.DefaultOrganization
	}
	if teamName == "" {
		teamName = constant.DefaultTeamName
	}
	return &GlobalConfig{
		Organization: organization,
		TeamName:     teamName,
	}
}

func NewGlobalConfigDefault() *GlobalConfig {
	return &GlobalConfig{
		Organization: constant.DefaultOrganization,
		TeamName:     constant.DefaultTeamName,
	}
}

func (p *GlobalConfig) IsComplete() bool {
	return p.Organization != "" && p.TeamName != ""
}
