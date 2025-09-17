package values

import "strings"

import shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"

const CATEGORY_PROJECT_DEFAULT_VALUE = "backend"

type CategoryProject struct {
	shared.BaseString
}

func NewCategoryProject(value string) (CategoryProject, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return NewDefaultCategoryProject(), nil
	}
	valueSafe := shared.MakeSafeForFileSystem(value)
	base, err := shared.NewBaseString(valueSafe, "CategoryProject")
	if err != nil {
		return CategoryProject{}, err
	}
	return CategoryProject{BaseString: base}, nil
}

func NewDefaultCategoryProject() CategoryProject {
	defaultCategory, _ := NewCategoryProject(CATEGORY_PROJECT_DEFAULT_VALUE)
	return defaultCategory
}

func (pn CategoryProject) Equals(other CategoryProject) bool {
	return pn.BaseString.Equals(other.BaseString)
}
