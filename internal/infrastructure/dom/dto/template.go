package dto

type TemplateDTO struct {
	RepositoryURL string `yaml:"repository_url"`
	Ref           string `yaml:"ref"`
}
