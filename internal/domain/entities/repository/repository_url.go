package repository

import (
	"errors"
	"net/url"
	"strings"
)

type RepositoryURL struct {
	value string
}

func NewRepositoryURL(value string) (RepositoryURL, error) {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return RepositoryURL{}, errors.New("RepositoryURL cannot be empty")
	}

	_, err := url.Parse(trimmedValue)
	if err != nil {
		return RepositoryURL{}, errors.New("RepositoryURL is not a valid URL")
	}

	if !strings.HasPrefix(trimmedValue, "http://") && !strings.HasPrefix(trimmedValue, "https://") {
		return RepositoryURL{}, errors.New("RepositoryURL must be HTTP or HTTPS")
	}

	return RepositoryURL{value: trimmedValue}, nil
}

func (r RepositoryURL) Value() string {
	return r.value
}

func (r RepositoryURL) String() string {
	return r.value
}

func (r RepositoryURL) Equals(other RepositoryURL) bool {
	return r.value == other.value
}
