package application

import (
	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/entities"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/ports"
)

type AppLogger struct {
	repo      ports.LoggerRepository
	presenter appPor.Presenter
}

func NewAppLogger(repo ports.LoggerRepository, presenter appPor.Presenter) appPor.Logger {
	return &AppLogger{
		repo:      repo,
		presenter: presenter,
	}
}

func (l *AppLogger) ShowLog() error {
	logger, err := l.repo.Find()
	if err != nil {
		return err
	}

	l.presenter.Header(&logger, logger.Revision())
	for _, step := range logger.Steps() {
		l.presenter.Step(step)
		for _, task := range step.Tasks() {
			l.presenter.Task(task, step)
		}
	}
	l.presenter.FinalSummary(&logger)

	return nil
}

func (l *AppLogger) StartExecution(contextData map[string]string, revision string) (*aggregates.Logger, error) {
	log := aggregates.NewLogger(contextData, revision)
	log.Start()

	if err := l.repo.Save(log); err != nil {
		return nil, err
	}

	if l.presenter != nil {
		l.presenter.Header(log, revision)
	}
	return log, nil
}

func (l *AppLogger) AddStep(logger *aggregates.Logger, stepName string) (*entities.StepRecord, error) {
	step, err := entities.NewStepRecord(stepName)
	if err != nil {
		return nil, err
	}

	if err := logger.AddStep(step); err != nil {
		return nil, err
	}

	if err := l.repo.Save(logger); err != nil {
		return nil, err
	}

	return step, nil
}

func (l *AppLogger) AddTaskToStep(logger *aggregates.Logger, stepName, taskName string) (*entities.TaskRecord, error) {
	step, err := logger.GetStep(stepName)
	if err != nil {
		return nil, err
	}

	task, err := entities.NewTaskRecord(taskName)
	if err != nil {
		return nil, err
	}

	step.AddTask(task)

	if err := l.repo.Save(logger); err != nil {
		return nil, err
	}

	return task, nil
}

func (l *AppLogger) MarkStepAsSuccessful(logger *aggregates.Logger, step *entities.StepRecord) error {
	step.MarkAsSuccess()
	if l.presenter != nil {
		l.presenter.Step(step)
	}
	return l.repo.Save(logger)
}

func (l *AppLogger) MarkStepAsFailed(logger *aggregates.Logger, step *entities.StepRecord, stepErr error) error {
	step.MarkAsFailure(stepErr)
	if l.presenter != nil {
		l.presenter.Step(step)
	}
	return l.repo.Save(logger)
}

func (l *AppLogger) MarkStepAsSkipped(logger *aggregates.Logger, step *entities.StepRecord) error {
	step.MarkAsSkipped()
	if l.presenter != nil {
		l.presenter.Step(step)
	}
	return l.repo.Save(logger)
}

func (l *AppLogger) MarkStepAsCached(logger *aggregates.Logger, step *entities.StepRecord, reason string) error {
	step.MarkAsCached(reason)
	if l.presenter != nil {
		l.presenter.Step(step)
	}
	return l.repo.Save(logger)
}

func (l *AppLogger) MarkStepAsRunning(logger *aggregates.Logger, step *entities.StepRecord) error {
	step.MarkAsRunning()
	if l.presenter != nil {
		l.presenter.Step(step)
	}
	return l.repo.Save(logger)
}

func (l *AppLogger) MarkTaskAsSuccessful(logger *aggregates.Logger, task *entities.TaskRecord, step *entities.StepRecord) error {
	task.MarkAsSuccess()
	if l.presenter != nil {
		l.presenter.Task(task, step)
	}
	return l.repo.Save(logger)
}

func (l *AppLogger) MarkTaskAsFailed(logger *aggregates.Logger, task *entities.TaskRecord, taskErr error, step *entities.StepRecord) error {
	task.MarkAsFailure(taskErr)
	if l.presenter != nil {
		l.presenter.Task(task, step)
	}
	return l.repo.Save(logger)
}

func (l *AppLogger) MarkTaskAsRunning(logger *aggregates.Logger, task *entities.TaskRecord, step *entities.StepRecord) error {
	task.MarkAsRunning()
	if l.presenter != nil {
		l.presenter.Task(task, step)
	}
	return l.repo.Save(logger)
}

func (l *AppLogger) SetTaskCommand(logger *aggregates.Logger, task *entities.TaskRecord, command string) error {
	task.SetCommand(command)
	return l.repo.Save(logger)
}

func (l *AppLogger) AddOutputToTask(logger *aggregates.Logger, task *entities.TaskRecord, outputLine string) error {
	task.AddOutput(outputLine)
	return l.repo.Save(logger)
}

func (l *AppLogger) FinishExecution(logger *aggregates.Logger) error {
	logger.RecalculateStatus()
	if l.presenter != nil {
		l.presenter.FinalSummary(logger)
	}
	return l.repo.Save(logger)
}
