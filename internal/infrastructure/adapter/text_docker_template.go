package adapter

import (
	"strings"
	"text/template"
)

type DockerTemplate interface {
	GetContent(pathTemplate string, params any) (string, error)
}

type TextDockerTemplate struct{}

func NewTextDockerTemplate() DockerTemplate {
	return &TextDockerTemplate{}
}

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
