package aggregates

import (
	"fmt"
	"time"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/entities"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/logger/vos"
)

type Logger struct {
	id        vos.LoggerID           `yaml:"id"`
	status    vos.Status             `yaml:"status"`
	startTime time.Time              `yaml:"start_time"`
	endTime   time.Time              `yaml:"end_time,omitempty"`
	steps     []*entities.StepRecord `yaml:"steps"`
	context   map[string]string      `yaml:"context"`
	stepIndex map[string]int
}

func NewLogger(id vos.LoggerID, context map[string]string) *Logger {
	return &Logger{
		id:        id,
		status:    vos.Pending,
		steps:     []*entities.StepRecord{},
		stepIndex: make(map[string]int),
		context:   context,
	}
}

func (e *Logger) RebuildIndex() {
	if e.stepIndex == nil {
		e.stepIndex = make(map[string]int)
	}
	for i, step := range e.steps {
		e.stepIndex[step.Name()] = i
	}
}

func (e *Logger) Start() {
	if e.status == vos.Pending {
		e.status = vos.Running
		e.startTime = time.Now()
	}
}

func (e *Logger) AddStep(step *entities.StepRecord) error {
	if _, exists := e.stepIndex[step.Name()]; exists {
		return fmt.Errorf("step with name '%s' already exists", step.Name())
	}
	e.steps = append(e.steps, step)
	e.stepIndex[step.Name()] = len(e.steps) - 1
	return nil
}

func (e *Logger) GetStep(name string) (*entities.StepRecord, error) {
	index, exists := e.stepIndex[name]
	if !exists {
		return nil, fmt.Errorf("step with name '%s' not found", name)
	}
	return e.steps[index], nil
}

func (e *Logger) ID() vos.LoggerID {
	return e.id
}

func (e *Logger) Context() map[string]string {
	return e.context
}

func (e *Logger) Steps() []*entities.StepRecord {
	return e.steps
}

func (e *Logger) Status() vos.Status {
	e.RecalculateStatus()
	return e.status
}

func (e *Logger) RecalculateStatus() {
	if e.status == vos.Success || e.status == vos.Failure || e.status == vos.Cached || e.status == vos.Skipped {
		return
	}

	hasFailure := false
	allFinished := true

	for _, step := range e.steps {
		stepStatus := step.Status()
		if stepStatus == vos.Failure {
			hasFailure = true
			break
		}
		if stepStatus != vos.Success {
			allFinished = false
		}
	}

	if hasFailure {
		if e.status == vos.Running {
			e.endTime = time.Now()
			e.status = vos.Failure
		}
		return
	}

	if allFinished {
		if e.status == vos.Running {
			e.status = vos.Success
			e.endTime = time.Now()
		}
		return
	}

	e.status = vos.Running
}
