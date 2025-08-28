package values

import (
	"strings"

	shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"
)

const REPOSITORY_VERSION_DEFAULT_VALUE = "v1.0.0"

type VersionRepository struct {
	shared.BaseString
}

func NewVersionRepository(value string) (VersionRepository, error) {
	value = strings.TrimSpace(value)
	
	if value == "" {
		return NewDefaultVersionRepository(), nil
	}
	
	valueSafe := shared.MakeSafeForFileSystem(value)
	base, err := shared.NewBaseString(valueSafe, "RepositoryVersion")
	if err != nil {
		return VersionRepository{}, err
	}
	return VersionRepository{BaseString: base}, nil
}

func NewDefaultVersionRepository() VersionRepository {
	return VersionRepository{BaseString: shared.NewBaseStringEmpty()}
}

func (tv VersionRepository) Equals(other VersionRepository) bool {
	return tv.BaseString.Equals(other.BaseString)
}
