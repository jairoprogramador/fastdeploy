package model

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type StepLog struct {
	StepName string
	Status   string
	Message  string
	Error    error
	Commands []string
}

type LogStore struct {
	Command  string
	LastStep StepLog
	Steps    []StepLog
}

func NewLogStore(command string) *LogStore {
	return &LogStore{
		Command: command,
		LastStep: StepLog{
			StepName: command,
		},
		Steps: []StepLog{},
	}
}

func (s *LogStore) StartStep(stepName string) {
	if s.LastStep.StepName != stepName && stepName != "" {
		s.Steps = append(s.Steps, s.LastStep)
		s.LastStep = StepLog{
			StepName: stepName,
			Status:   StatusSuccess,
			Message:  "",
			Error:    nil,
			Commands: []string{},
		}
	}
}

func (s *LogStore) AddCommand(command string) {
	s.LastStep.Commands = append(s.LastStep.Commands, command)
}

func (s *LogStore) AddMessage(message string) {
	if message != "" {
		s.LastStep.Message = message
		s.LastStep.Status = StatusSuccess
	}
}

func (s *LogStore) AddError(err error) {
	if err != nil {
		s.LastStep.Error = err
		s.LastStep.Status = StatusError
	}
}

func (s *LogStore) FinishSteps() {
	if s.LastStep.StepName != "" {
		s.Steps = append(s.Steps, s.LastStep)
		s.LastStep = StepLog{}
	}
}

func (s *LogStore) GetLogs() []StepLog {
	s.FinishSteps()
	return s.Steps
}

func (s *LogStore) HasErrors() bool {
	s.FinishSteps()
	for _, step := range s.Steps {
		if step.Status == StatusError || step.Error != nil {
			return true
		}
	}
	return false
}

func (s *LogStore) GetError() error {
	s.FinishSteps()
	for _, step := range s.Steps {
		if step.Status == StatusError || step.Error != nil {
			return step.Error
		}
	}
	return nil
}



