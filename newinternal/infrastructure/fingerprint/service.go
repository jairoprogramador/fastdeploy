package fingerprint

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/vos"
	gitignore "github.com/sabhiram/go-gitignore"
)

// Service implementa la interfaz services.FingerprintService.
type Service struct{}

// NewService crea una nueva instancia del servicio de fingerprint.
func NewService() *Service {
	return &Service{}
}

// CalculateCodeFingerprint calcula un hash del contenido de un directorio, respetando un archivo .fdignore.
func (s *Service) CalculateCodeFingerprint(_ context.Context, projectPath string, ignoreFiles []string) (vos.Fingerprint, error) {
	ignoreMatcher, err := gitignore.CompileIgnoreFile(filepath.Join(projectPath, ".fdignore"))
	if err != nil && !os.IsNotExist(err) {
		return vos.Fingerprint{}, fmt.Errorf("error al leer el archivo .fdignore: %w", err)
	}
	if err != nil { // El archivo no existe, creamos un matcher vacío.
		ignoreMatcher = &gitignore.GitIgnore{}
	}

	return s.calculateDirectoryHash(projectPath, ignoreMatcher)
}

// CalculateEnvironmentFingerprint calcula un hash del contenido de un directorio sin ignorar archivos.
func (s *Service) CalculateEnvironmentFingerprint(_ context.Context, environmentPath string) (vos.Fingerprint, error) {
	return s.calculateDirectoryHash(environmentPath, nil)
}

// calculateDirectoryHash es la lógica central para hashear un directorio.
func (s *Service) calculateDirectoryHash(rootPath string, ignorer *gitignore.GitIgnore) (vos.Fingerprint, error) {
	fileHashes := make(map[string]string)

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil // Solo procesamos archivos
		}

		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}

		if ignorer != nil && ignorer.MatchesPath(relPath) {
			return nil
		}

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

func (s *Service) hashFile(filePath string) (string, error) {
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

func (s *Service) combineHashes(fileHashes map[string]string) (string, error) {
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
