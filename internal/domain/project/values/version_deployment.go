package values

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"
)

var versionRegex = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)(-[a-zA-Z0-9.-]+)?$`)

type VersionDeployment struct {
	shared.BaseString
}

func NewVersionDeployment(value string) (VersionDeployment, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return NewDefaultVersionDeployment(), nil
	}
	
	err := validateVersion(value)
	if err != nil {
		return VersionDeployment{}, err
	}

	base, err := shared.NewBaseString(value, "DeploymentVersion")
	if err != nil {
		return VersionDeployment{}, err
	}
	return VersionDeployment{BaseString: base}, nil
}

func NewDefaultVersionDeployment() VersionDeployment {
	defaultVersion, _ := NewVersionDeployment("v1.0.0")
	return defaultVersion
}

func (d VersionDeployment) Increment() VersionDeployment {
	matches := versionRegex.FindStringSubmatch(d.Value())

	if len(matches) < 4 {
		return NewDefaultVersionDeployment()
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	patch++

	newVersion := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	version, _ := NewVersionDeployment(newVersion)
	return version
}

func (d VersionDeployment) Equals(other VersionDeployment) bool {
	return d.BaseString.Equals(other.BaseString)
}

func (d VersionDeployment) IsValid() bool {
	return validateVersion(d.Value()) == nil
}

func validateVersion(value string) error {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return errors.New("DeploymentVersion cannot be empty")
	}

	if !versionRegex.MatchString(trimmedValue) {
		return errors.New("DeploymentVersion must follow format vX.Y.Z-suffix")
	}

	return nil
}
