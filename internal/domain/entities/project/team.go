package project

import "github.com/jairoprogramador/fastdeploy/internal/domain/entities/common"

type Team struct {
	common.StringValueObject
}

func NewTeam(value string) (Team, error) {
	base, err := common.NewStringValueObject(value, "Team")
	if err != nil {
		return Team{}, err
	}
	return Team{StringValueObject: base}, nil
}
