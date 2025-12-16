package vos

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"
)

type TemplateRepository struct {
	url string
	ref string
}

func NewTemplateRepository(repoURL, ref string) (TemplateRepository, error) {
	if repoURL == "" {
		return TemplateRepository{}, errors.New("la URL del repositorio de plantillas no puede estar vacía")
	}

	repoURLConverted := repoURL
	if strings.HasPrefix(repoURL, "git@") && !strings.HasPrefix(repoURL, "ssh://") {
		repoURLConverted = "ssh://" + strings.Replace(repoURL, ":", "/", 1)
	}

	parsedURL, err := url.Parse(repoURLConverted)
	if err != nil {
		return TemplateRepository{}, errors.New("la URL del repositorio de plantillas no es válida")
	}

	// Si después de parsear no hay esquema, es probable que sea una ruta local o una URL inválida.
	if parsedURL.Scheme == "" {
		return TemplateRepository{}, errors.New("la URL del repositorio debe tener un esquema (ej: https, ssh)")
	}

	if ref == "" {
		ref = "main"
	}

	return TemplateRepository{
		url: repoURL,
		ref: ref,
	}, nil
}

func (t TemplateRepository) URL() string {
	return t.url
}

func (t TemplateRepository) Ref() string {
	return t.ref
}

func (t TemplateRepository) DirName() string {
	base := filepath.Base(t.url)
	return strings.TrimSuffix(base, ".git")
}
