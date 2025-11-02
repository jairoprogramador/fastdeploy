package ports

import (
	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/entities"
)

type Logger interface {
	ShowLog(pathProject string) error
	Start(namesParams appDto.NamesParams, contextData map[string]string, revision string) (*aggregates.Logger, error)
	AddStep(namesParams appDto.NamesParams, logger *aggregates.Logger, stepName string) (*entities.StepRecord, error)
	AddTaskToStep(namesParams appDto.NamesParams, logger *aggregates.Logger, stepName, taskName string) (*entities.TaskRecord, error)
	MarkStepAsSuccessful(namesParams appDto.NamesParams, logger *aggregates.Logger, step *entities.StepRecord) error
	MarkStepAsFailed(namesParams appDto.NamesParams, logger *aggregates.Logger, stepRecord *entities.StepRecord, stepErr error) error
	MarkStepAsSkipped(namesParams appDto.NamesParams, logger *aggregates.Logger, stepRecord *entities.StepRecord) error
	MarkStepAsCached(namesParams appDto.NamesParams, logger *aggregates.Logger, stepRecord *entities.StepRecord, reason string) error
	MarkStepAsRunning(namesParams appDto.NamesParams, logger *aggregates.Logger, stepRecord *entities.StepRecord) error
	MarkTaskAsSuccessful(namesParams appDto.NamesParams, logger *aggregates.Logger, taskRecord *entities.TaskRecord, stepRecord *entities.StepRecord) error
	MarkTaskAsFailed(namesParams appDto.NamesParams, logger *aggregates.Logger, taskRecord *entities.TaskRecord, taskErr error, stepRecord *entities.StepRecord) error
	MarkTaskAsRunning(namesParams appDto.NamesParams, logger *aggregates.Logger, taskRecord *entities.TaskRecord, stepRecord *entities.StepRecord) error
	SetTaskCommand(namesParams appDto.NamesParams, logger *aggregates.Logger, taskRecord *entities.TaskRecord, command string) error
	AddOutputToTask(namesParams appDto.NamesParams, logger *aggregates.Logger, taskRecord *entities.TaskRecord, outputLine string) error
	FinishExecution(namesParams appDto.NamesParams, logger *aggregates.Logger) error
}
