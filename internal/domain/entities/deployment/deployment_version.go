package deployment

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	VersionPattern = `^v(\d+)\.(\d+)\.(\d+)(-[a-zA-Z0-9.-]+)?$`
)

type DeploymentVersion struct {
	value string
}

func NewDeploymentVersion(value string) (DeploymentVersion, error) {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return DeploymentVersion{}, errors.New("DeploymentVersion cannot be empty")
	}

	pattern := regexp.MustCompile(VersionPattern)
	if !pattern.MatchString(trimmedValue) {
		return DeploymentVersion{}, errors.New("DeploymentVersion must follow format vX.Y.Z-suffix")
	}

	return DeploymentVersion{value: trimmedValue}, nil
}

func CreateInitialVersion() DeploymentVersion {
	return DeploymentVersion{value: "v1.0.0"}
}

func (d DeploymentVersion) Increment() DeploymentVersion {
	pattern := regexp.MustCompile(VersionPattern)
	matches := pattern.FindStringSubmatch(d.value)

	if len(matches) < 4 {
		return CreateInitialVersion()
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	patch++

	newVersion := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	return DeploymentVersion{value: newVersion}
}

func (d DeploymentVersion) Value() string {
	return d.value
}

func (d DeploymentVersion) String() string {
	return d.value
}

func (d DeploymentVersion) Equals(other DeploymentVersion) bool {
	return d.value == other.value
}
