package values

import (
	shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"
)

const ORGANIZATION_DEFAULT_VALUE = "fastdeploy"

type NameOrganization struct {
	shared.BaseString
}

func NewNameOrganization(value string) (NameOrganization, error) {
	valueSafe := shared.MakeSafeForFileSystem(value)
	base, err := shared.NewBaseString(valueSafe, "Organization")
	if err != nil {
		return NameOrganization{}, err
	} 
	return NameOrganization{BaseString: base}, nil
}

func NewDefaultNameOrganization() NameOrganization {
	defaultOrganization, _ := NewNameOrganization(ORGANIZATION_DEFAULT_VALUE)
	return defaultOrganization
}

func (o NameOrganization) Equals(other NameOrganization) bool {
	return o.BaseString.Equals(other.BaseString)
}
