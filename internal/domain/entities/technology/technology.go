package technology

import (
	"fmt"
	"path/filepath"
)

type Technology struct {
	name    TechnologyName
	version TechnologyVersion
}

func NewTechnology(name TechnologyName, version TechnologyVersion) Technology {
	return Technology{
		name:    name,
		version: version,
	}
}

func (t Technology) GetName() TechnologyName {
	return t.name
}

func (t Technology) GetVersion() TechnologyVersion {
	return t.version
}

func (t Technology) GetPath() string {
	return filepath.Join(t.name.Value(), t.version.Value())
}

func (t Technology) GetFullPath(step string) string {
	return filepath.Join(step, t.GetPath())
}

func (t Technology) IsValid() bool {
	return t.name.Value() != "" && t.version.Value() != ""
}

func (t Technology) Equals(other Technology) bool {
	return t.name.StringValueObject.Equals(other.name.StringValueObject) &&
		t.version.StringValueObject.Equals(other.version.StringValueObject)
}

func (t Technology) String() string {
	return fmt.Sprintf("%s-%s", t.name.Value(), t.version.Value())
}
