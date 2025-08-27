package values

import (
	"strings"

	shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"
)
type NameTechnology struct {
	shared.BaseString
}

func NewNameTechnology(value string) (NameTechnology, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return NewDefaultNameTechnology(), nil
	}
	
	valueSafe := shared.MakeSafeForFileSystem(value)
	base, err := shared.NewBaseString(valueSafe, "TechnologyName")
	if err != nil {
		return NameTechnology{}, err
	}
	return NameTechnology{BaseString: base}, nil
}

func NewDefaultNameTechnology() NameTechnology {
	return NameTechnology{BaseString: shared.NewBaseStringEmpty()}
}

func (tn NameTechnology) Equals(other NameTechnology) bool {
	return tn.BaseString.Equals(other.BaseString)
}
