package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/jairoprogramador/fastdeploy/internal/domain/project/port"
)

type HashIdentifier struct{}

func NewHashIdentifier() port.Identifier {
	return &HashIdentifier{}
}

func (s *HashIdentifier) Generate(projectName string, organizationName string) string {
	nameProject := strings.ToLower(strings.TrimSpace(projectName))
	nameOrganization := strings.ToLower(strings.TrimSpace(organizationName))

	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%s-%s-%d", nameProject, nameOrganization, timestamp)

	hash := sha256.Sum256([]byte(data))
	hashString := hex.EncodeToString(hash[:])

	return hashString
}
