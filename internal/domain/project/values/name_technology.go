package values

import shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"

const TECHNOLOGY_NAME_DEFAULT_VALUE = "springboot"

type NameTechnology struct {
	shared.BaseString
}

func NewNameTechnology(value string) (NameTechnology, error) {
	valueSafe := shared.MakeSafeForFileSystem(value)
	base, err := shared.NewBaseString(valueSafe, "TechnologyName")
	if err != nil {
		return NameTechnology{}, err
	}
	return NameTechnology{BaseString: base}, nil
}

func NewDefaultNameTechnology() NameTechnology {
	defaultTechnologyName, _ := NewNameTechnology(TECHNOLOGY_NAME_DEFAULT_VALUE)
	return defaultTechnologyName
}

func (tn NameTechnology) Equals(other NameTechnology) bool {
	return tn.BaseString.Equals(other.BaseString)
}
