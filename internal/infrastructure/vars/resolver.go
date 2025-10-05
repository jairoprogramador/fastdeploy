package vars

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	deploymentvos "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/services"
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
)

// varRegex es la expresión regular compilada para encontrar placeholders de variables como ${var.nombre_variable}.
var varRegex = regexp.MustCompile(`\$\{var\.([a-zA-Z0-9_.-]+)\}`)

// Resolver implementa la interfaz services.VariableResolver.
// Este adaptador es responsable de toda la lógica de interpolación y extracción de variables.
type Resolver struct{}

// NewResolver crea una nueva instancia del Resolver.
func NewResolver() services.VariableResolver {
	return &Resolver{}
}

// Interpolate reemplaza los placeholders de variables en una cadena con sus valores reales.
func (r *Resolver) Interpolate(template string, variables map[string]vos.Variable) (string, error) {
	//fmt.Println("----------------Interpolate----------------", variables)
	var firstErr error
	result := varRegex.ReplaceAllStringFunc(template, func(match string) string {
		key := varRegex.FindStringSubmatch(match)[1]
		variable, ok := variables[key]
		if !ok {
			if firstErr == nil {
				firstErr = fmt.Errorf("variable '%s' no encontrada en el mapa de variables", key)
			}
			return match // Si no se encuentra, dejar el placeholder original
		}
		return variable.Value()
	})

	if firstErr != nil {
		return "", firstErr
	}
	return result, nil
}

// ProcessTemplateFile procesa un único archivo o recorre un directorio para procesar todos los archivos que contiene.
func (r *Resolver) ProcessTemplate(path string, variables map[string]vos.Variable) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error al obtener información de la ruta '%s': %w", path, err)
	}

	if !info.IsDir() {
		// El path es un archivo, procesarlo directamente.
		return r.processSingleFile(path, variables)
	}

	// El path es un directorio, recorrerlo recursivamente.
	return filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Procesar solo archivos, no directorios.
		if !d.IsDir() {
			err := r.processSingleFile(currentPath, variables)
			if err != nil {
				// Envolver el error con el path del archivo que falló.
				return fmt.Errorf("fallo al procesar el archivo '%s': %w", currentPath, err)
			}
		}
		return nil
	})
}

// processSingleFile contiene la lógica para procesar un único archivo de plantilla "in-place".
func (r *Resolver) processSingleFile(filePath string, variables map[string]vos.Variable) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("error al obtener información del archivo de plantilla '%s': %w", filePath, err)
	}
	fileMode := info.Mode()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo de plantilla '%s': %w", filePath, err)
	}

	interpolatedContent, err := r.Interpolate(string(content), variables)
	if err != nil {
		return fmt.Errorf("error al interpolar variables en el archivo '%s': %w", filePath, err)
	}

	// Sobrescribir el archivo original con el contenido interpolado y los permisos originales.
	err = os.WriteFile(filePath, []byte(interpolatedContent), fileMode)
	if err != nil {
		return fmt.Errorf("error al escribir las modificaciones en el archivo de plantilla '%s': %w", filePath, err)
	}

	return nil
}

// ExtractVariable utiliza una sonda (expresión regular) para encontrar y/o extraer una nueva variable de un texto.
func (r *Resolver) ExtractVariable(probe deploymentvos.OutputProbe, text string) (vos.Variable, bool, error) {
	re, err := regexp.Compile(probe.Probe())
	if err != nil {
		return vos.Variable{}, false, fmt.Errorf("expresión regular de la sonda no es válida '%s': %w", probe.Probe(), err)
	}

	matches := re.FindStringSubmatch(text)
	if matches == nil {
		return vos.Variable{}, false, fmt.Errorf("no se encontró ninguna coincidencia para la sonda '%s'", probe.Probe())
	}

	if probe.Name() == "" {
		return vos.Variable{}, true, nil
	}

	// Si hay un nombre, se espera que la regex tenga al menos un grupo de captura para el valor.
	if len(matches) < 2 {
		return vos.Variable{}, false, fmt.Errorf(
			"la sonda '%s' coincidió, pero no se encontró un grupo de captura para extraer el valor",
			probe.Name(),
		)
	}

	// El valor extraído es el contenido del primer grupo de captura.
	value := matches[1]
	variable, err := vos.NewVariable(probe.Name(), value)
	if err != nil {
		return vos.Variable{}, false, fmt.Errorf("error al crear la variable extraída '%s': %w", probe.Name(), err)
	}

	return variable, true, nil
}
