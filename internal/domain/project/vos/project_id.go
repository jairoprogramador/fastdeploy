package vos

import (
	"crypto/sha256"
	"fmt"
)

type ProjectID struct {
	value string
}

func NewProjectID(id string) ProjectID {
	return ProjectID{value: id}
}

func GenerateProjectID(name, organization, team string) ProjectID {
	data := fmt.Sprintf("%s-%s-%s", name, organization, team)
	hash := sha256.Sum256([]byte(data))
	return ProjectID{value: fmt.Sprintf("%x", hash)}
}

func (p ProjectID) String() string {
	return p.value
}

func (p ProjectID) Equals(other ProjectID) bool {
	return p.value == other.value
}
