package workspace

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"github.com/jairoprogramador/fastdeploy/internal/application/ports"
)

// Manager implementa la interfaz ports.WorkspaceManager.
// Es responsable de preparar los directorios de trabajo para la ejecución de cada paso.
type Manager struct {
	projectsBasePath string
	respositoriesBasePath string
}

// NewManager crea una nueva instancia del WorkspaceManager.
func NewManager(projectsBasePath string, respositoriesBasePath string) ports.WorkspaceManager {
	return &Manager{
		projectsBasePath: projectsBasePath,
		respositoriesBasePath: respositoriesBasePath,
	}
}

// PrepareStepWorkspace crea un directorio de trabajo limpio y copia los archivos de la plantilla del paso.
func (m *Manager) PrepareStepWorkspace(
	projectName string,
	environment string,
	stepName string,
	repositoryName string,
) (string, error) {
	// 1. Encontrar el directorio de origen del paso en la plantilla.
	sourceDir, err := m.findStepSourceDir(repositoryName, stepName)
	if err != nil {
		// Si no hay un directorio de plantilla para el paso, no es un error.
		// Simplemente creamos un directorio de trabajo vacío.
		if os.IsNotExist(err) {
			sourceDir = ""
		} else {
			return "", err
		}
	}

	// 2. Construir la ruta de destino.
	destPath := filepath.Join(m.projectsBasePath, projectName, environment, stepName)

	// 3. Limpiar y recrear el directorio de destino para asegurar un estado prístino.
	/* if err := os.RemoveAll(destPath); err != nil {
		return "", fmt.Errorf("error al limpiar el workspace del paso anterior: %w", err)
	} */
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return "", fmt.Errorf("error al crear el workspace del paso: %w", err)
	}

	// 4. Copiar los archivos de la plantilla si existe un directorio de origen.
	if sourceDir != "" {
		if err := copyDir(sourceDir, destPath); err != nil {
			return "", fmt.Errorf("error al copiar los archivos de la plantilla al workspace: %w", err)
		}
	}

	return destPath, nil
}

// findStepSourceDir busca en el directorio "steps" el subdirectorio que corresponde a un nombre de paso.
// Ej: para stepName "test", busca un directorio como "01-test".
func (m *Manager) findStepSourceDir(repositoryName, stepName string) (string, error) {
	stepsRoot := filepath.Join(m.respositoriesBasePath, repositoryName, "steps")
	entries, err := os.ReadDir(stepsRoot)
	if err != nil {
		return "", err
	}

	regex := regexp.MustCompile(fmt.Sprintf(`^\d+-%s$`, regexp.QuoteMeta(stepName)))
	for _, entry := range entries {
		if entry.IsDir() && regex.MatchString(entry.Name()) {
			return filepath.Join(stepsRoot, entry.Name()), nil
		}
	}

	return "", os.ErrNotExist // Devuelve un error específico si no se encuentra.
}

// copyDir copia el contenido de un directorio a otro de forma recursiva.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath, info)
	})
}

// copyFile copia un único archivo.
func copyFile(src, dst string, info os.FileInfo) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
