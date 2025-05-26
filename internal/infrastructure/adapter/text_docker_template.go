package adapter

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model/logger"
	"strings"
	"text/template"
)

type DockerTemplate interface {
	GetContent(pathTemplate string, params any) (string, error)
}

type TextDockerTemplate struct {
	fileLogger *logger.FileLogger
}

func NewTextDockerTemplate(fileLogger *logger.FileLogger) DockerTemplate {
	return &TextDockerTemplate{
		fileLogger: fileLogger,
	}
}

func (t *TextDockerTemplate) GetContent(pathTemplate string, params any) (string, error) {
	templateFile, err := template.ParseFiles(pathTemplate)
	if err != nil {
		t.logError(err)
		return "", err
	}

	var result strings.Builder
	err = templateFile.Execute(&result, params)
	if err != nil {
		t.logError(err)
		return "", err
	}

	return result.String(), nil
}

func (t *TextDockerTemplate) logError(err error) {
	if err != nil {
		t.fileLogger.Error(err)
	}
}
