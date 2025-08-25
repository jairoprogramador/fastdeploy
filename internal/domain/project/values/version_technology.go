package values

import (
	"strings"

	shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"
)

const TECHNOLOGY_VERSION_DEFAULT_VALUE = "3.5.4"

type VersionTechnology struct {
	shared.BaseString
}

func NewVersionTechnology(value string) (VersionTechnology, error) {
	value = strings.TrimSpace(value)
	
	if value == "" {
		return NewDefaultVersionTechnology(), nil
	}
	
	valueSafe := shared.MakeSafeForFileSystem(value)
	base, err := shared.NewBaseString(valueSafe, "TechnologyVersion")
	if err != nil {
		return VersionTechnology{}, err
	}
	return VersionTechnology{BaseString: base}, nil
}

func NewDefaultVersionTechnology() VersionTechnology {
	defaultVersion, _ := NewVersionTechnology(TECHNOLOGY_VERSION_DEFAULT_VALUE)
	return defaultVersion
}

func (tv VersionTechnology) Equals(other VersionTechnology) bool {
	return tv.BaseString.Equals(other.BaseString)
}
