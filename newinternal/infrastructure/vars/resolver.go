package vars

import (
	"fmt"
	"os"
	"regexp"

	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

// varRegex es la expresión regular compilada para encontrar placeholders de variables como ${var.nombre_variable}.
var varRegex = regexp.MustCompile(`\$\{var\.([a-zA-Z0-9_.-]+)\}`)

// Resolver implementa la interfaz services.VariableResolver.
// Este adaptador es responsable de toda la lógica de interpolación y extracción de variables.
type Resolver struct{}

// NewResolver crea una nueva instancia del Resolver.
func NewResolver() *Resolver {
	return &Resolver{}
}

// Interpolate reemplaza los placeholders de variables en una cadena con sus valores reales.
func (r *Resolver) Interpolate(template string, variables map[string]vos.Variable) (string, error) {
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

// ProcessTemplateFile lee un archivo, interpola sus variables y escribe el resultado en un nuevo archivo.
func (r *Resolver) ProcessTemplateFile(srcPath, destPath string, variables map[string]vos.Variable) error {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo de plantilla '%s': %w", srcPath, err)
	}

	interpolatedContent, err := r.Interpolate(string(content), variables)
	if err != nil {
		return fmt.Errorf("error al interpolar variables en el archivo '%s': %w", srcPath, err)
	}

	err = os.WriteFile(destPath, []byte(interpolatedContent), 0644)
	if err != nil {
		return fmt.Errorf("error al escribir el archivo de plantilla procesado en '%s': %w", destPath, err)
	}

	return nil
}

// ExtractVariable utiliza una sonda (expresión regular) para encontrar y/o extraer una nueva variable de un texto.
func (r *Resolver) ExtractVariable(probe deploymentvos.OutputProbe, text string) (vos.Variable, bool, error) {
	re, err := regexp.Compile(probe.Probe())
	if err != nil {
		// Este error no debería ocurrir si el VO se construye correctamente, pero es una buena salvaguarda.
		return vos.Variable{}, false, fmt.Errorf("expresión regular de la sonda no es válida '%s': %w", probe.Probe(), err)
	}

	matches := re.FindStringSubmatch(text)
	if matches == nil {
		return vos.Variable{}, false, fmt.Errorf("la sonda '%s' no coincidió", re.String())
	}

	// Si el nombre de la sonda está vacío, solo nos importaba si había una coincidencia.
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
		// Error al crear el VO (e.g., el valor extraído estaba vacío).
		return vos.Variable{}, false, fmt.Errorf("error al crear la variable extraída '%s': %w", probe.Name(), err)
	}

	return variable, true, nil
}
