package template

import (
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"strings"
	"text/template"
)

type DockerTemplatePort interface {
	GetContent(pathTemplate string, params any) (string, error)
}

type templateAdapter struct {
	fileLogger *logger.FileLogger
}

func NewTemplateAdapter(fileLogger *logger.FileLogger) DockerTemplatePort {
	return &templateAdapter{
		fileLogger: fileLogger,
	}
}

func (t *templateAdapter) GetContent(pathTemplate string, params any) (string, error) {
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

func (t *templateAdapter) logError(err error) {
	if err != nil {
		t.fileLogger.Error(err)
	}
}
