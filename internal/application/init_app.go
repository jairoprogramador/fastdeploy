package application

import (
	"deploy/internal/domain/model"
	"deploy/internal/domain/service"
)

func InitApp(projectService service.ProjectServiceInterface, logStore *model.LogStore) *model.LogStore {
	message, err := projectService.Initialize()
	if err != nil {
		logStore.AddError(err)
	} else {
		logStore.AddMessage(message)
	}
	logStore.FinishSteps()
	return logStore
}
