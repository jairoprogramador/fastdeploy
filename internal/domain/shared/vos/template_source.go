package vos

import (
	"errors"
	"strings"
)

type TemplateSource struct {
	url string
	ref     string
}

func NewTemplateSource(url, ref string) (TemplateSource, error) {
	if url == "" {
		return TemplateSource{}, errors.New("la url del template no puede estar vacía")
	}
	if ref == "" {
		return TemplateSource{}, errors.New("la referencia del template no puede estar vacía")
	}
	return TemplateSource{url: url, ref: ref}, nil
}

func (ts TemplateSource) Url() string {
	return ts.url
}

func (ts TemplateSource) Ref() string {
	return ts.ref
}

func (ts TemplateSource) NameTemplate() string {
	safePath := strings.Split(ts.url, "/")
	lastPart := safePath[len(safePath)-1]
	return strings.TrimSuffix(lastPart, ".git")
}