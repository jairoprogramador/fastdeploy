package command

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/jairoprogramador/fastdeploy/internal/domain/command/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
)

type TemplateFile struct{}

func NewTemplateFile() port.TemplatePort {
	return &TemplateFile{}
}

func (e *TemplateFile) Process(pathTemplate string, processor port.LineProcessor, context *values.ContextValue) error {
	infoFile, err := os.Stat(pathTemplate)
	if err != nil {
		return err
	}

	if infoFile.IsDir() {
		err := filepath.Walk(pathTemplate, func(currentPath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if fileInfo.IsDir() {
				return nil
			}

			pathTemplateFile, err := filepath.Rel(pathTemplate, currentPath)
			if err != nil {
				return err
			}
			return e.processFile(pathTemplateFile, processor, context)
		})
		if err != nil {
			return err
		}
	} else {
		return e.processFile(pathTemplate, processor, context)
	}

	return nil
}

func (e *TemplateFile) processFile(
	pathFileTemplate string,
	processor port.LineProcessor,
	context *values.ContextValue) error {

	fileTemplate, err := os.Open(pathFileTemplate)
	if err != nil {
		return err
	}
	defer fileTemplate.Close()

	var linesFileTemplate []string
	scannerTemplate := bufio.NewScanner(fileTemplate)

	for scannerTemplate.Scan() {
		lineFileTemplate := scannerTemplate.Text()
		processedLine := processor(lineFileTemplate, context)
		linesFileTemplate = append(linesFileTemplate, processedLine)
	}

	if err := scannerTemplate.Err(); err != nil {
		return err
	}

	output := strings.Join(linesFileTemplate, "\n")
	err = os.WriteFile(pathFileTemplate, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}
