package project

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/cli/presenter"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/service"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

func StartDeploy(projectService service.ProjectService, fileLogger *logger.FileLogger, ctx context.Context) result.DomainResult {
	presenter.ShowMessage("environment: local")
	spinner := presenter.ShowMessageSpinner("implementation in progress")

	result := projectService.Start(ctx)
	fileLogger.WriteToFile()

	defer spinner.Stop()
	return result
}
