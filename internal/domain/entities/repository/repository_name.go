package repository

import (
	"strings"

	"github.com/jairoprogramador/fastdeploy/internal/domain/entities/common"
)

type RepositoryName struct {
	common.StringValueObject
}

func NewRepositoryName(value string) (RepositoryName, error) {
	base, err := common.NewStringValueObject(value, "RepositoryName")
	if err != nil {
		return RepositoryName{}, err
	}
	return RepositoryName{StringValueObject: base}, nil
}

func ExtractFromURL(url RepositoryURL) (RepositoryName, error) {
	urlPath := url.Value()
	parts := strings.Split(urlPath, "/")

	lastPart := parts[len(parts)-1]
	name := strings.TrimSuffix(lastPart, ".git")

	return NewRepositoryName(name)
}
