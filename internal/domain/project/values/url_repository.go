package values

import (
	"errors"
	"regexp"
	"strings"

	shared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/values"
)

var urlRegex = regexp.MustCompile(`^(https?:\/\/[^\s]+|git@[^\s:]+:[^\s]+)$`)

const GIT_URL_DEFAULT_VALUE = "https://github.com/jairoprogramador/mydeploy.git"

type UrlRepository struct {
	shared.BaseString
}

func NewUrlRepository(value string) (UrlRepository, error) {
	err := validateUrl(value)
	if err != nil {
		return UrlRepository{}, err
	}
	
	base, err := shared.NewBaseString(value, "RepositoryURL")
	if err != nil {
		return UrlRepository{}, err
	}
	return UrlRepository{BaseString: base}, nil
}

func NewDefaultUrlRepository() UrlRepository {
	defaultRepo, _ := NewUrlRepository(GIT_URL_DEFAULT_VALUE)
	return defaultRepo
}

func (r UrlRepository) Equals(other UrlRepository) bool {
	return r.BaseString.Equals(other.BaseString)
}

func (r UrlRepository) IsValid() bool {
	return validateUrl(r.Value()) == nil
}

func (r UrlRepository) ExtractNameRepository() string {
	urlPath := r.Value()
	parts := strings.Split(urlPath, "/")

	lastPart := parts[len(parts)-1]
	name := strings.TrimSuffix(lastPart, ".git")

	return name
}

func validateUrl(value string) error {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return errors.New("repository.url cannot be empty")
	}

	if !urlRegex.MatchString(trimmedValue) {
		return errors.New("repository.url must be HTTP, HTTPS or GIT")
	}

	return nil
}
