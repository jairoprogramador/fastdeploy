package reader

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type CommandDefinition struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Cmd         string `yaml:"cmd"`
}

type CommandConfig struct {
	Commands []CommandDefinition `yaml:"commands"`
}

func (c *CommandConfig) UnmarshalYAML(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func (c *CommandConfig) ReadFileYAML(yamlFilePath string) error {
	data, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo de comandos YAML en %s: %w", yamlFilePath, err)
	}

	if err := c.UnmarshalYAML(data); err != nil {
		return fmt.Errorf("error al deserializar el archivo de comandos: %w", err)
	}

	return nil
}
