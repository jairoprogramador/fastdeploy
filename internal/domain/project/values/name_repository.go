package values

import (
	shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"
)

type NameRepository struct {
	shared.BaseString
}

func NewNameRepository(value string) (NameRepository, error) {
	base, err := shared.NewBaseString(value, "RepositoryName")
	if err != nil {
		return NameRepository{}, err
	}
	return NameRepository{BaseString: base}, nil
}

func (rn NameRepository) Equals(other NameRepository) bool {
	return rn.BaseString.Equals(other.BaseString)
}
