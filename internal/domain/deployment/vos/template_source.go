package vos

import "errors"

// TemplateSource representa el origen de una plantilla de despliegue.
// Es un Objeto de Valor que da identidad a un DeploymentTemplate.
type TemplateSource struct {
	repoURL string
	ref     string // Commit hash, tag, o branch
}

// NewTemplateSource crea un nuevo y validado Objeto de Valor TemplateSource.
func NewTemplateSource(repoURL, ref string) (TemplateSource, error) {
	if repoURL == "" {
		return TemplateSource{}, errors.New("la URL del repositorio no puede estar vacía")
	}
	if ref == "" {
		return TemplateSource{}, errors.New("la referencia (commit/tag/branch) no puede estar vacía")
	}
	return TemplateSource{repoURL: repoURL, ref: ref}, nil
}

// RepoURL devuelve la URL del repositorio de la plantilla.
func (ts TemplateSource) RepoURL() string {
	return ts.repoURL
}

// Ref devuelve la referencia (commit/tag/branch) de la plantilla.
func (ts TemplateSource) Ref() string {
	return ts.ref
}
