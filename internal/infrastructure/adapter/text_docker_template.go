package adapter

import (
	"deploy/internal/domain/port"
	"strings"
	"text/template"
)

type TextDockerTemplate struct{}

// NewTextDockerTemplate creates a new instance of DockerTemplate
func NewTextDockerTemplate() port.DockerTemplate {
	return &TextDockerTemplate{}
}

// GetContentTemplate parses a template file and executes it with the provided parameters
func (t *TextDockerTemplate) GetContent(pathTemplate string, params any) (string, error) {
	templateFile, err := template.ParseFiles(pathTemplate)
	if err == nil {
		var result strings.Builder
		err = templateFile.Execute(&result, params)
		if err == nil {
			return result.String(), nil
		}
		return "", err
	}
	return "", err
}
