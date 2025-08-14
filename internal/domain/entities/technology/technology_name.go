package technology

import "github.com/jairoprogramador/fastdeploy/internal/domain/entities/common"

type TechnologyName struct {
	common.StringValueObject
}

func NewTechnologyName(value string) (TechnologyName, error) {
	base, err := common.NewStringValueObject(value, "TechnologyName")
	if err != nil {
		return TechnologyName{}, err
	}
	return TechnologyName{StringValueObject: base}, nil
}
