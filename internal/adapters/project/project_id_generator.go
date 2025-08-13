package project

import (
	"crypto/sha1"
	"fmt"
	"time"
)

type ProjectIDGenerator interface {
	GenerateID(projectName string) (string, error)
}

type ProjectIDGeneratorImpl struct{}

func NewProjectIDGenerator() ProjectIDGenerator {
	return &ProjectIDGeneratorImpl{}
}

func (spg *ProjectIDGeneratorImpl) GenerateID(projectName string) (string, error) {
	timestamp := time.Now().String()
	data := []byte(projectName + timestamp)

	hash := sha1.New()
	_, err := hash.Write(data)
	if err != nil {
		return "", fmt.Errorf("error al generar el hash SHA-1: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
