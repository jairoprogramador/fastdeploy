package values

import (
	shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"
)

type Identifier struct {
	shared.BaseString
}

func NewIdentifier(value string) (Identifier, error) {
	base, err := shared.NewBaseString(value, "ProjectID")
	if err != nil {
		return Identifier{}, err
	}
	return Identifier{BaseString: base}, nil
}

func (p Identifier) Equals(other Identifier) bool {
	return p.BaseString.Equals(other.BaseString)
}
