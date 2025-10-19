package services

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	staSer "github.com/jairoprogramador/fastdeploy/internal/domain/state/services"
	staVos "github.com/jairoprogramador/fastdeploy/internal/domain/state/vos"

	gitignore "github.com/sabhiram/go-gitignore"
)

type FingerprintService struct {
	environment    string
}

func NewFingerprintService(
	environment string,
) (staSer.FingerprintService, error) {

	if environment == "" {
		return nil, fmt.Errorf("environment is required")
	}

	return &FingerprintService{
		environment:    environment,
	}, nil
}

func (s *FingerprintService) GenerateFromSource(sourcePath string) (staVos.Fingerprint, error) {
	pathGitIgnore := filepath.Join(sourcePath, ".gitignore")
	lines := []string{".git", ".gitignore"}
	return s.generateHash(sourcePath, pathGitIgnore, lines)
}

func (s *FingerprintService) GenerateFromStepDefinition(templatePath, stepName string) (staVos.Fingerprint, error) {
	pathGitIgnore := filepath.Join(templatePath, ".gitignore")
	lines := []string{".git", ".gitignore"}

	stepsPathToIgnore, err := s.getStepsPathToIgnore(templatePath, stepName)
	if err != nil {
		return staVos.Fingerprint{}, err
	}
	lines = append(lines, stepsPathToIgnore...)

	varsPathToIgnore, err := s.getVarsPathToIgnore(templatePath, stepName)
	if err != nil {
		return staVos.Fingerprint{}, err
	}
	lines = append(lines, varsPathToIgnore...)

	return s.generateHash(templatePath, pathGitIgnore, lines)
}

func (s *FingerprintService) GenerateFromStepVariables(vars map[string]string) (staVos.Fingerprint, error) {
	hashString, err := s.generateHashStringForMap(vars)
	if err != nil {
		return staVos.Fingerprint{}, err
	}
	return staVos.NewFingerprint(hashString)
}

func (s *FingerprintService) getStepsPathToIgnore(templatePath, stepName string) ([]string, error) {
	stepsPath := filepath.Join(templatePath, "steps")

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

func (s *FingerprintService) getVarsPathToIgnore(templatePath,stepName string) ([]string, error) {
	variablesPath := filepath.Join(templatePath, "variables")

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
			variablesStepFileToIgnore, err := s.getVariablesStepFileToIgnore(templatePath, stepName)
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

func (s *FingerprintService) getVariablesStepFileToIgnore(templatePath, stepName string) ([]string, error) {
	var variablesStepFileToIgnore []string

	pathStepVariables := filepath.Join(templatePath, "variables", s.environment)

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

func (s *FingerprintService) generateHash(
	sourcePath string,
	pathGitIgnore string,
	ignoreLines []string) (staVos.Fingerprint, error){

	ignoreMatcher, err := gitignore.CompileIgnoreFileAndLines(pathGitIgnore, ignoreLines...)
	if err != nil && !os.IsNotExist(err) {
		return staVos.Fingerprint{}, fmt.Errorf("error al leer el archivo .gitignore: %w", err)
	}
	if err != nil {
		ignoreMatcher = &gitignore.GitIgnore{}
	}
	return s.generateHashFromSource(sourcePath, ignoreMatcher)
}

func (s *FingerprintService) generateHashFromSource(sourcePath string, ignorer *gitignore.GitIgnore) (staVos.Fingerprint, error) {
	fileHashes := make(map[string]string)

	err := filepath.WalkDir(sourcePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourcePath, path)
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
		return staVos.Fingerprint{}, err
	}

	finalHash, err := s.generateHashStringForMap(fileHashes)
	if err != nil {
		return staVos.Fingerprint{}, err
	}

	return staVos.NewFingerprint(finalHash)
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

func (s *FingerprintService) generateHashStringForMap(data map[string]string) (string, error) {
	if len(data) == 0 {
		return "d41d8cd98f00b204e9800998ecf8427e", nil
	}

	paths := make([]string, 0, len(data))
	for path := range data {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	masterHash := sha256.New()
	for _, path := range paths {
		line := fmt.Sprintf("%s:%s\n", path, data[path])
		if _, err := masterHash.Write([]byte(line)); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%x", masterHash.Sum(nil)), nil
}
