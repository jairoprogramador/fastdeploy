package fingerprint

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/services"
	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/vos"
	gitignore "github.com/sabhiram/go-gitignore"
)

// Service implementa la interfaz services.FingerprintService.
type FingerprintService struct{}

// NewService crea una nueva instancia del servicio de fingerprint.
func NewFingerprintService() services.FingerprintService {
	return &FingerprintService{}
}

// CalculateCodeFingerprint calcula un hash del contenido de un directorio, respetando un archivo .fdignore.
func (s *FingerprintService) CalculateCodeFingerprint(_ context.Context, pathProject string) (vos.Fingerprint, error) {
	pathGitIgnore := filepath.Join(pathProject, ".gitignore")
	lines := []string{".git", ".gitignore"}

	ignoreMatcher, err := gitignore.CompileIgnoreFileAndLines(pathGitIgnore, lines...)
	if err != nil && !os.IsNotExist(err) {
		return vos.Fingerprint{}, fmt.Errorf("error al leer el archivo .gitignore: %w", err)
	}
	if err != nil { // El archivo no existe, creamos un matcher vacío.
		ignoreMatcher = &gitignore.GitIgnore{}
	}

	return s.calculateDirectoryHash(pathProject, ignoreMatcher)
}

// CalculateEnvironmentFingerprint calcula un hash del contenido de un directorio sin ignorar archivos.
func (s *FingerprintService) CalculateEnvironmentFingerprint(_ context.Context, stepName string, pathRepository string) (vos.Fingerprint, error) {

	pathGitIgnore := filepath.Join(pathRepository, ".gitignore")
	lines := []string{".git", ".gitignore"}

	// Agregar directorios de steps a ignorar (excepto el stepName actual)
	stepsToIgnore, err := s.getStepsToIgnore(stepName, pathRepository)
	if err != nil {
		return vos.Fingerprint{}, fmt.Errorf("error al obtener directorios de steps a ignorar: %w", err)
	}
	lines = append(lines, stepsToIgnore...)

	ignoreMatcher, err := gitignore.CompileIgnoreFileAndLines(pathGitIgnore, lines...)
	if err != nil && !os.IsNotExist(err) {
		return vos.Fingerprint{}, fmt.Errorf("error al leer el archivo .gitignore: %w", err)
	}
	if err != nil { // El archivo no existe, creamos un matcher vacío.
		ignoreMatcher = &gitignore.GitIgnore{}
	}
	return s.calculateDirectoryHash(pathRepository, ignoreMatcher)
}

// getStepsToIgnore obtiene los directorios de steps que deben ser ignorados
// excepto el stepName actual. Los directorios siguen el patrón: "01-test", "02-supply", etc.
func (s *FingerprintService) getStepsToIgnore(stepName string, pathRepository string) ([]string, error) {
	stepsPath := filepath.Join(pathRepository, "steps")

	// Verificar si el directorio steps existe
	if _, err := os.Stat(stepsPath); os.IsNotExist(err) {
		return []string{}, nil // No hay directorio steps, no hay nada que ignorar
	}

	entries, err := os.ReadDir(stepsPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer el directorio steps: %w", err)
	}

	// Crear regex para el patrón: dígitos + guión + stepName
	// Escapar caracteres especiales en stepName para regex
	escapedStepName := regexp.QuoteMeta(stepName)
	pattern := fmt.Sprintf(`^\d+-%s$`, escapedStepName)
	stepRegex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("error al compilar regex para stepName '%s': %w", stepName, err)
	}

	var stepsToIgnore []string

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()

		// Si el directorio NO coincide con el patrón del stepName actual, lo ignoramos
		if !stepRegex.MatchString(dirName) {
			stepsToIgnore = append(stepsToIgnore, filepath.Join("steps", dirName, "**.*"))
		}
	}
	return stepsToIgnore, nil
}

// calculateDirectoryHash es la lógica central para hashear un directorio.
func (s *FingerprintService) calculateDirectoryHash(rootPath string, ignorer *gitignore.GitIgnore) (vos.Fingerprint, error) {
	fileHashes := make(map[string]string)

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}

		if ignorer != nil && ignorer.MatchesPath(relPath) {
			return nil
		}
		//fmt.Println("-----------path-----------", path)
		hash, err := s.hashFile(path)
		if err != nil {
			return fmt.Errorf("fallo al hashear el archivo %s: %w", path, err)
		}
		fileHashes[relPath] = hash
		return nil
	})

	if err != nil {
		return vos.Fingerprint{}, err
	}

	// Combinar todos los hashes de archivos en un hash final y estable.
	finalHash, err := s.combineHashes(fileHashes)
	if err != nil {
		return vos.Fingerprint{}, err
	}

	return vos.NewFingerprint(finalHash)
}

func (s *FingerprintService) hashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (s *FingerprintService) combineHashes(fileHashes map[string]string) (string, error) {
	if len(fileHashes) == 0 {
		return "d41d8cd98f00b204e9800998ecf8427e", nil // Hash SHA-1 de una cadena vacía, un valor constante.
	}

	// Ordenar las rutas de los archivos es CRÍTICO para un hash estable.
	paths := make([]string, 0, len(fileHashes))
	for path := range fileHashes {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	masterHash := sha256.New()
	for _, path := range paths {
		line := fmt.Sprintf("%s:%s\n", path, fileHashes[path])
		if _, err := masterHash.Write([]byte(line)); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%x", masterHash.Sum(nil)), nil
}
