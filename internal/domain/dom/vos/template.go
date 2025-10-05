package vos

import "errors"

// Template agrupa los datos que definen un template. Es un Objeto de Valor.
type Template struct {
	repository_url string
	ref            string
}

func NewTemplate(url, ref string) (*Template, error) {
	if url == "" {
		return nil, errors.New("la URL del repositorio no puede estar vacío")
	}
	if ref == "" {
		return nil, errors.New("la referencia no puede estar vacía")
	}
	return &Template{
		repository_url: url,
		ref:            ref,
	}, nil
}

// Getters para todos los campos...
func (t *Template) RepositoryURL() string { return t.repository_url }
func (t *Template) Ref() string { return t.ref }