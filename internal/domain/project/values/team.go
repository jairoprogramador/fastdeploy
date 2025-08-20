package values

import (
	shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"
)

const TEAM_DEFAULT_VALUE = "itachi"

type Team struct {
	shared.BaseString
}

func NewTeam(value string) (Team, error) {
	valueSafe := shared.MakeSafeForFileSystem(value)
	base, err := shared.NewBaseString(valueSafe, "Team")
	if err != nil {
		return Team{}, err
	}
	return Team{BaseString: base}, nil
}

func NewDefaultTeam() Team {
	defaultTeam, _ := NewTeam(TEAM_DEFAULT_VALUE)
	return defaultTeam
}

func (t Team) Equals(other Team) bool {
	return t.BaseString.Equals(other.BaseString)
}
