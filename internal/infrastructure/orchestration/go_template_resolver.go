package orchestration

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/services"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/vos"
)

var varRegex = regexp.MustCompile(`\$\{var\.([a-zA-Z0-9_.-]+)\}`)

type GoTemplateResolver struct{}

func NewGoTemplateResolver() services.TemplateResolver {
	return &GoTemplateResolver{}
}

func (r *GoTemplateResolver) ResolveTemplate(template string, variables map[string]vos.Output) (string, error) {
	var firstErr error
	result := varRegex.ReplaceAllStringFunc(template, func(match string) string {
		key := varRegex.FindStringSubmatch(match)[1]
		variable, ok := variables[key]
		if !ok {
			if firstErr == nil {
				firstErr = fmt.Errorf("variable '%s' no encontrada en el mapa de variables", key)
			}
			return match
		}
		return variable.Value()
	})

	if firstErr != nil {
		return "", firstErr
	}
	return result, nil
}

func (r *GoTemplateResolver) ResolvePath(path string, variables map[string]vos.Output) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error al obtener información de la ruta '%s': %w", path, err)
	}

	if !info.IsDir() {
		return r.resolverTemplateFile(path, variables)
	}

	return filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			err := r.resolverTemplateFile(currentPath, variables)
			if err != nil {
				return fmt.Errorf("fallo al procesar el archivo '%s': %w", currentPath, err)
			}
		}
		return nil
	})
}

func (r *GoTemplateResolver) resolverTemplateFile(filePath string, variables map[string]vos.Output) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("error al obtener información del archivo de plantilla '%s': %w", filePath, err)
	}
	fileMode := info.Mode()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo de plantilla '%s': %w", filePath, err)
	}

	interpolatedContent, err := r.ResolveTemplate(string(content), variables)
	if err != nil {
		return fmt.Errorf("error al interpolar variables en el archivo '%s': %w", filePath, err)
	}

	err = os.WriteFile(filePath, []byte(interpolatedContent), fileMode)
	if err != nil {
		return fmt.Errorf("error al escribir las modificaciones en el archivo de plantilla '%s': %w", filePath, err)
	}

	return nil
}

func (r *GoTemplateResolver) ResolveOutput(output vos.Output, record string) (vos.Output, bool, error) {
	re, err := regexp.Compile(output.Value())
	if err != nil {
		return vos.Output{}, false, fmt.Errorf("expresión regular de la sonda no es válida '%s': %w", output.Value(), err)
	}

	matches := re.FindStringSubmatch(record)
	if matches == nil {
		return vos.Output{}, false, fmt.Errorf("no se encontró ninguna coincidencia para la sonda '%s'", output.Value())
	}

	if output.Name() == "" {
		return vos.Output{}, true, nil
	}

	if len(matches) < 2 {
		return vos.Output{}, false, fmt.Errorf(
			"la sonda '%s' coincidió, pero no se encontró un grupo de captura para extraer el valor",
			output.Name(),
		)
	}

	value := matches[1]
	outputExtracted, err := vos.NewOutputFromNameAndValue(output.Name(), value)
	if err != nil {
		return vos.Output{}, false, err
	}

	return outputExtracted, true, nil
}
