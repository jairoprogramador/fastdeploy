package values

type Technology struct {
	name    NameTechnology
	version VersionTechnology
}

func NewTechnology(name NameTechnology, version VersionTechnology) Technology {
	return Technology{
		name:    name,
		version: version,
	}
}

func NewDefaultTechnology() Technology {
	return NewTechnology(NewDefaultNameTechnology(), NewDefaultVersionTechnology())
}

func (t Technology) GetName() NameTechnology {
	return t.name
}

func (t Technology) GetVersion() VersionTechnology {
	return t.version
}

/* func (t Technology) GetPath() string {
	return filepath.Join(t.name.BaseString.Value(), t.version.BaseString.Value())
}

func (t Technology) GetFullPath(step string) string {
	return filepath.Join(step, t.GetPath())
} */

func (t Technology) Equals(other Technology) bool {
	return t.name.Equals(other.name) &&
		t.version.Equals(other.version)
}

