package application

import (
	"deploy/internal/domain/service"
)

func InitApp(projectService service.ProjectService) error {
	_, err := projectService.Initialize()
	return err
}
