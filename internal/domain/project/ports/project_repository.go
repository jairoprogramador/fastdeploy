package ports

import (
	"context"
)

type ProjectConfigDTO struct {
	ID           string
	Name         string
	Organization string
	Team         string
	Description  string
	Version      string
	TemplateURL  string
	TemplateRef  string
}

type ProjectRepository interface {
	Load(ctx context.Context, pathFile string) (*ProjectConfigDTO, error)
	Save(ctx context.Context, pathFile string, data *ProjectConfigDTO) error
}
