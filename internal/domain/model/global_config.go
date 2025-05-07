package model

type GlobalConfig struct {
	Organization string `yaml:"organization"`
	TeamName     string `yaml:"teamName"`
}

func NewGlobalConfig(organization, teamName string) *GlobalConfig {
	if organization == "" {
		organization = DefaultOrganization
	}
	if teamName == "" {
		teamName = DefaultTeamName
	}
	return &GlobalConfig{
		Organization: organization,
		TeamName:     teamName,
	}
}

func NewGlobalConfigDefault() *GlobalConfig {
	return &GlobalConfig{
		Organization: DefaultOrganization,
		TeamName:     DefaultTeamName,
	}
}

func (p *GlobalConfig) IsComplete() bool {
	return p.Organization != "" && p.TeamName != ""
}
