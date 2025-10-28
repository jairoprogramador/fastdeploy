package ports

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/entities"
)

type Presenter interface {
	Header(log *aggregates.Logger, revision string)
	Step(step *entities.StepRecord)
	Task(task *entities.TaskRecord, step *entities.StepRecord)
	FinalSummary(log *aggregates.Logger)
}
