package dto

// Variable representa una única variable leída desde un archivo YAML.
type Variable struct {
	Name  string      `yaml:"name"`
	Value interface{} `yaml:"value"`
}
