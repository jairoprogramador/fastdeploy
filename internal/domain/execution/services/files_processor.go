package services

import (
	"fmt"
	"os"
	"path/filepath"

	execvos "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
	sharedVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared/vos"
)

type TemplateProcessor struct {
	interpolator  *Interpolator
	workspaceRoot string
	backups       map[string][]byte
}

func NewTemplateProcessor(workspaceRoot string, interpolator *Interpolator) *TemplateProcessor {
	return &TemplateProcessor{
		interpolator:  interpolator,
		workspaceRoot: workspaceRoot,
		backups:       make(map[string][]byte),
	}
}

func (tp *TemplateProcessor) Process(templates []*sharedVos.FileTemplate, vars execvos.VariableSet) error {
	for _, tpl := range templates {
		absPath := filepath.Join(tp.workspaceRoot, tpl.Path)

		if _, exists := tp.backups[absPath]; !exists {
			originalContent, err := os.ReadFile(absPath)
			if err != nil {
				return fmt.Errorf("no se pudo leer el archivo de plantilla original %s: %w", absPath, err)
			}
			tp.backups[absPath] = originalContent
		}

		interpolatedContent, err := tp.interpolator.Interpolate(string(tp.backups[absPath]), vars)
		if err != nil {
			return fmt.Errorf("no se pudo interpolar la plantilla %s: %w", absPath, err)
		}

		if err := os.WriteFile(absPath, []byte(interpolatedContent), 0644); err != nil {
			return fmt.Errorf("no se pudo escribir el archivo de plantilla interpolado %s: %w", absPath, err)
		}
	}
	return nil
}

func (tp *TemplateProcessor) Restore() error {
	var firstErr error
	for path, originalContent := range tp.backups {
		if err := os.WriteFile(path, originalContent, 0644); err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("no se pudo restaurar el archivo de plantilla %s: %w", path, err)
			}
		}
	}
	return firstErr
}
