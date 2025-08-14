package technology

import "github.com/jairoprogramador/fastdeploy/internal/domain/entities/common"

type TechnologyVersion struct {
	common.StringValueObject
}

func NewTechnologyVersion(value string) (TechnologyVersion, error) {
	base, err := common.NewStringValueObject(value, "TechnologyVersion")
	if err != nil {
		return TechnologyVersion{}, err
	}
	return TechnologyVersion{StringValueObject: base}, nil
}
