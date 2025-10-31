package ports

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/entities"
)

type Logger interface {
	ShowLog() error
	StartExecution(contextData map[string]string, revision string) (*aggregates.Logger, error)
	AddStep(logger *aggregates.Logger, stepName string) (*entities.StepRecord, error)
	AddTaskToStep(logger *aggregates.Logger, stepName, taskName string) (*entities.TaskRecord, error)
	MarkStepAsSuccessful(logger *aggregates.Logger, step *entities.StepRecord) error
	MarkStepAsFailed(logger *aggregates.Logger, step *entities.StepRecord, stepErr error) error
	MarkStepAsSkipped(logger *aggregates.Logger, step *entities.StepRecord) error
	MarkStepAsCached(logger *aggregates.Logger, step *entities.StepRecord, reason string) error
	MarkStepAsRunning(logger *aggregates.Logger, step *entities.StepRecord) error
	MarkTaskAsSuccessful(logger *aggregates.Logger, task *entities.TaskRecord, step *entities.StepRecord) error
	MarkTaskAsFailed(logger *aggregates.Logger, task *entities.TaskRecord, taskErr error, step *entities.StepRecord) error
	MarkTaskAsRunning(logger *aggregates.Logger, task *entities.TaskRecord, step *entities.StepRecord) error
	SetTaskCommand(logger *aggregates.Logger, task *entities.TaskRecord, command string) error
	AddOutputToTask(logger *aggregates.Logger, task *entities.TaskRecord, outputLine string) error
	FinishExecution(logger *aggregates.Logger) error
}