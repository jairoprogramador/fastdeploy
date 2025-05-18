package application

import (
	"deploy/internal/domain/model"
)

var (
	projectModel *model.Project
)

func Initialize() *model.LogStore {
	logStore := model.NewLogStore("initialize project")
	message, err := getProjectService().Initialize()
	if err != nil {
		logStore.AddError(err)
	} else {
		logStore.AddMessage(message)
	}
	logStore.FinishSteps()
	return logStore
}

func IsInitialize() error {
	var err error
	projectModel, err = getProjectService().Load()
	return err
}
