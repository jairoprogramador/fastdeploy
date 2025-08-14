package project

import "github.com/jairoprogramador/fastdeploy/internal/domain/entities/common"

type ProjectName struct {
	common.StringValueObject
}

func NewProjectName(value string) (ProjectName, error) {
	base, err := common.NewStringValueObject(value, "ProjectName")
	if err != nil {
		return ProjectName{}, err
	}
	return ProjectName{StringValueObject: base}, nil
}
