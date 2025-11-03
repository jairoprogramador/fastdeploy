package vos

import "strings"

const (
	DefaultUrl = "https://github.com/jairoprogramador/mydeploy.git"
	DefaultRef = "main"
)

type Template struct {
	url string
	ref string
}

func NewTemplate(url, ref string) Template {
	if url == "" {
		url = DefaultUrl
	}
	if ref == "" {
		ref = DefaultRef
	}
	return Template{url: url, ref: ref}
}

func (t Template) URL() string { return t.url }
func (t Template) Ref() string { return t.ref }

func (t Template) NameTemplate() string {
	safePath := strings.Split(t.url, "/")
	lastPart := safePath[len(safePath)-1]
	return strings.TrimSuffix(lastPart, ".git")
}