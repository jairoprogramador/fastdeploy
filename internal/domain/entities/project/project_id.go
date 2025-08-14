package project

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type ProjectID struct {
	value string
}

func NewProjectID(value string) (ProjectID, error) {
	if value == "" {
		return ProjectID{}, errors.New("ProjectID cannot be empty")
	}
	return ProjectID{value: value}, nil
}

func GenerateProjectID(projectName, organization string) ProjectID {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%s-%s-%d", projectName, organization, timestamp)

	hash := sha256.Sum256([]byte(data))
	hashString := hex.EncodeToString(hash[:])

	return ProjectID{value: hashString[:16]}
}

func (p ProjectID) Value() string {
	return p.value
}

func (p ProjectID) String() string {
	return p.value
}

func (p ProjectID) Equals(other ProjectID) bool {
	return p.value == other.value
}
