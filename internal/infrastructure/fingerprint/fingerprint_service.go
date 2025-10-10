package fingerprint

import (
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

type FingerprintService struct {
	pathProjectApp string
	pathRepository string
	environment    string
}

func NewFingerprintService(
	pathProjectApp string,
	pathRepositoryRootFastDeploy string,
	repositoryName string,
	environment string,
) (services.FingerprintService, error) {

	if pathProjectApp == "" {
		return nil, fmt.Errorf("path project root fast deploy is required")
	}
	if pathRepositoryRootFastDeploy == "" {
		return nil, fmt.Errorf("path repository root fast deploy is required")
	}
	if repositoryName == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if environment == "" {
		return nil, fmt.Errorf("environment is required")
	}

	pathRepository := filepath.Join(pathRepositoryRootFastDeploy, repositoryName)

	return &FingerprintService{
		pathProjectApp: pathProjectApp,
		pathRepository: pathRepository,
		environment:    environment,
	}, nil
}

func (s *FingerprintService) CalculateCodeFingerprint() (vos.Fingerprint, error) {
	pathGitIgnore := filepath.Join(s.pathProjectApp, ".gitignore")
	lines := []string{".git", ".gitignore"}

	ignoreMatcher, err := gitignore.CompileIgnoreFileAndLines(pathGitIgnore, lines...)
	if err != nil && !os.IsNotExist(err) {
		return vos.Fingerprint{}, fmt.Errorf("error al leer el archivo .gitignore: %w", err)
	}
	if err != nil {
		ignoreMatcher = &gitignore.GitIgnore{}
	}

	return s.calculateDirectoryHash(s.pathProjectApp, ignoreMatcher)
}

func (s *FingerprintService) CalculateStepFingerprint(stepName string) (vos.Fingerprint, error) {
	pathGitIgnore := filepath.Join(s.pathRepository, ".gitignore")
	lines := []string{".git", ".gitignore"}

	stepsToIgnore, err := s.getStepsToIgnore(stepName)
	if err != nil {
		return vos.Fingerprint{}, err
	}
	lines = append(lines, stepsToIgnore...)

	variablesToIgnore, err := s.getVariablesFileToIgnore(stepName)
	if err != nil {
		return vos.Fingerprint{}, err
	}
	lines = append(lines, variablesToIgnore...)

	ignoreMatcher, err := gitignore.CompileIgnoreFileAndLines(pathGitIgnore, lines...)
	if err != nil && !os.IsNotExist(err) {
		return vos.Fingerprint{}, fmt.Errorf("error al leer el archivo .gitignore: %w", err)
	}
	if err != nil {
		ignoreMatcher = &gitignore.GitIgnore{}
	}
	return s.calculateDirectoryHash(s.pathRepository, ignoreMatcher)
}

func (s *FingerprintService) getStepsToIgnore(stepName string) ([]string, error) {
	stepsPath := filepath.Join(s.pathRepository, "steps")

	if _, err := os.Stat(stepsPath); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(stepsPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer el directorio steps: %w", err)
	}

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

		if !stepRegex.MatchString(dirName) {
			stepsToIgnore = append(stepsToIgnore, filepath.Join("steps", dirName, "**.*"))
		}
	}
	return stepsToIgnore, nil
}

func (s *FingerprintService) getVariablesFileToIgnore(stepName string) ([]string, error) {
	variablesPath := filepath.Join(s.pathRepository, "variables")

	if _, err := os.Stat(variablesPath); os.IsNotExist(err) {
		return []string{}, nil
	}

	entriesVariables, err := os.ReadDir(variablesPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer el directorio variables: %w", err)
	}

	var variablesToIgnore []string

	for _, entryVariable := range entriesVariables {
		entryVariableName := entryVariable.Name()

		if entryVariableName == s.environment {
			variablesStepFileToIgnore, err := s.getVariablesStepFileToIgnore(stepName)
			if err != nil {
				return nil, fmt.Errorf("error al obtener directorios de variables del step: %w", err)
			}
			variablesToIgnore = append(variablesToIgnore, variablesStepFileToIgnore...)
			continue
		}

		if !entryVariable.IsDir() {
			variablesToIgnore = append(variablesToIgnore, filepath.Join("variables", entryVariableName))
			continue
		}

		if entryVariableName != s.environment {
			variablesToIgnore = append(variablesToIgnore, filepath.Join("variables", entryVariableName, "**.*"))
		}
	}
	return variablesToIgnore, nil
}

func (s *FingerprintService) getVariablesStepFileToIgnore(stepName string) ([]string, error) {
	var variablesStepFileToIgnore []string

	pathStepVariables := filepath.Join(s.pathRepository, "variables", s.environment)

	entriesStepVariables, err := os.ReadDir(pathStepVariables)
	if err != nil {
		return nil, fmt.Errorf("error al leer el directorio variables del step: %w", err)
	}

	for _, entryStepVariables := range entriesStepVariables {
		fileStepName := fmt.Sprintf("%s.yaml", stepName)

		if !entryStepVariables.IsDir() && entryStepVariables.Name() == fileStepName {
			continue
		}

		if !entryStepVariables.IsDir() {
			variablesStepFileToIgnore = append(
				variablesStepFileToIgnore,
				filepath.Join("variables", s.environment, entryStepVariables.Name()))
		} else {
			variablesStepFileToIgnore = append(
				variablesStepFileToIgnore,
				filepath.Join("variables", s.environment, entryStepVariables.Name(), "**.*"))
		}
	}
	return variablesStepFileToIgnore, nil
}

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
