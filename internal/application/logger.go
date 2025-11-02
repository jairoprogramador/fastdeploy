package application

import (
	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"

	proPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/ports"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/entities"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/ports"
)

type AppLogger struct {
	loggerRepository ports.LoggerRepository
	configRepository proPor.ConfigRepository
	presenter  appPor.Presenter
}

func NewAppLogger(
	loggerRepository ports.LoggerRepository,
	configRepository proPor.ConfigRepository,
	presenter appPor.Presenter) appPor.Logger {
	return &AppLogger{
		loggerRepository: loggerRepository,
		configRepository: configRepository,
		presenter:  presenter,
	}
}

func (l *AppLogger) ShowLog(pathProject string) error {

	configProject, err := l.configRepository.Load(pathProject)
	if err != nil {
		return err
	}

	namesParams := appDto.NewNamesParams(configProject.Project().Name(), configProject.Template().NameTemplate())

	logger, err := l.loggerRepository.Find(namesParams)
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

func (l *AppLogger) Start(namesParams appDto.NamesParams, contextData map[string]string, revision string) (*aggregates.Logger, error) {
	log := aggregates.NewLogger(contextData, revision)
	log.Start()

	if err := l.loggerRepository.Save(namesParams, log); err != nil {
		return nil, err
	}

	if l.presenter != nil {
		l.presenter.Header(log, revision)
	}
	return log, nil
}

func (l *AppLogger) AddStep(namesParams appDto.NamesParams, logger *aggregates.Logger, stepName string) (*entities.StepRecord, error) {
	step, err := entities.NewStepRecord(stepName)
	if err != nil {
		return nil, err
	}

	if err := logger.AddStep(step); err != nil {
		return nil, err
	}

	if err := l.loggerRepository.Save(namesParams, logger); err != nil {
		return nil, err
	}

	return step, nil
}

func (l *AppLogger) AddTaskToStep(namesParams appDto.NamesParams, logger *aggregates.Logger, stepName, taskName string) (*entities.TaskRecord, error) {
	step, err := logger.GetStep(stepName)
	if err != nil {
		return nil, err
	}

	task, err := entities.NewTaskRecord(taskName)
	if err != nil {
		return nil, err
	}

	step.AddTask(task)

	if err := l.loggerRepository.Save(namesParams, logger); err != nil {
		return nil, err
	}

	return task, nil
}

func (l *AppLogger) MarkStepAsSuccessful(namesParams appDto.NamesParams, logger *aggregates.Logger, step *entities.StepRecord) error {
	step.MarkAsSuccess()
	if l.presenter != nil {
		l.presenter.Step(step)
	}
	return l.loggerRepository.Save(namesParams, logger)
}

func (l *AppLogger) MarkStepAsFailed(namesParams appDto.NamesParams, logger *aggregates.Logger, step *entities.StepRecord, stepErr error) error {
	step.MarkAsFailure(stepErr)
	if l.presenter != nil {
		l.presenter.Step(step)
	}
	return l.loggerRepository.Save(namesParams, logger)
}

func (l *AppLogger) MarkStepAsSkipped(namesParams appDto.NamesParams, logger *aggregates.Logger, step *entities.StepRecord) error {
	step.MarkAsSkipped()
	if l.presenter != nil {
		l.presenter.Step(step)
	}
	return l.loggerRepository.Save(namesParams, logger)
}

func (l *AppLogger) MarkStepAsCached(namesParams appDto.NamesParams, logger *aggregates.Logger, step *entities.StepRecord, reason string) error {
	step.MarkAsCached(reason)
	if l.presenter != nil {
		l.presenter.Step(step)
	}
	return l.loggerRepository.Save(namesParams, logger)
}

func (l *AppLogger) MarkStepAsRunning(namesParams appDto.NamesParams, logger *aggregates.Logger, step *entities.StepRecord) error {
	step.MarkAsRunning()
	if l.presenter != nil {
		l.presenter.Step(step)
	}
	return l.loggerRepository.Save(namesParams, logger)
}

func (l *AppLogger) MarkTaskAsSuccessful(namesParams appDto.NamesParams, logger *aggregates.Logger, task *entities.TaskRecord, step *entities.StepRecord) error {
	task.MarkAsSuccess()
	if l.presenter != nil {
		l.presenter.Task(task, step)
	}
	return l.loggerRepository.Save(namesParams, logger)
}

func (l *AppLogger) MarkTaskAsFailed(namesParams appDto.NamesParams, logger *aggregates.Logger, task *entities.TaskRecord, taskErr error, step *entities.StepRecord) error {
	task.MarkAsFailure(taskErr)
	if l.presenter != nil {
		l.presenter.Task(task, step)
	}
	return l.loggerRepository.Save(namesParams, logger)
}

func (l *AppLogger) MarkTaskAsRunning(namesParams appDto.NamesParams, logger *aggregates.Logger, task *entities.TaskRecord, step *entities.StepRecord) error {
	task.MarkAsRunning()
	if l.presenter != nil {
		l.presenter.Task(task, step)
	}
	return l.loggerRepository.Save(namesParams, logger)
}

func (l *AppLogger) SetTaskCommand(namesParams appDto.NamesParams, logger *aggregates.Logger, task *entities.TaskRecord, command string) error {
	task.SetCommand(command)
	return l.loggerRepository.Save(namesParams, logger)
}

func (l *AppLogger) AddOutputToTask(namesParams appDto.NamesParams, logger *aggregates.Logger, task *entities.TaskRecord, outputLine string) error {
	task.AddOutput(outputLine)
	return l.loggerRepository.Save(namesParams, logger)
}

func (l *AppLogger) FinishExecution(namesParams appDto.NamesParams, logger *aggregates.Logger) error {
	logger.RecalculateStatus()
	if l.presenter != nil {
		l.presenter.FinalSummary(logger)
	}
	return l.loggerRepository.Save(namesParams, logger)
}
