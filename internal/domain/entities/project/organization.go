package project

import "github.com/jairoprogramador/fastdeploy/internal/domain/entities/common"

type Organization struct {
	common.StringValueObject
}

func NewOrganization(value string) (Organization, error) {
	base, err := common.NewStringValueObject(value, "Organization")
	if err != nil {
		return Organization{}, err
	}
	return Organization{StringValueObject: base}, nil
}
