package entity

type Environment struct {
	Name string
}

func NewEnvironment(name string) Environment {
	return Environment{Name: name}
}

func (e Environment) GetName() string {
	return e.Name
}

func (e Environment) Equals(other Environment) bool {
	return e.Name == other.Name
}
