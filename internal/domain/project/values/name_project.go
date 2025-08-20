package values

import shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"

type NameProject struct {
	shared.BaseString
}

func NewNameProject(value string) (NameProject, error) {
	valueSafe := shared.MakeSafeForFileSystem(value)
	base, err := shared.NewBaseString(valueSafe, "ProjectName")
	if err != nil {
		return NameProject{}, err
	}
	return NameProject{BaseString: base}, nil
}

func (pn NameProject) Equals(other NameProject) bool {
	return pn.BaseString.Equals(other.BaseString)
}
