package model

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type StepLog struct {
	StepName     string
	Status   string
	Message   string
	Error    error
	Commands []string
}

type LogStore struct {
	Command  string
	LastStep StepLog
	Steps    []StepLog
}

var (
	instanceLogStore *LogStore
)

func NewLogStore(command string) *LogStore {
	instanceLogStore = &LogStore{
		Command: command,
		LastStep: StepLog{
			StepName: command,
		},
		Steps:    []StepLog{},
	}
	return instanceLogStore
}

func GetLogStore() *LogStore {
	if instanceLogStore == nil {
		instanceLogStore = NewLogStore("not defined")
	}
	return instanceLogStore
}

func (s *LogStore) StartStep(stepName string) {
	if s.LastStep.StepName != stepName && stepName != "" {
		s.Steps = append(s.Steps, s.LastStep)
		s.LastStep = StepLog{
			StepName: stepName,
			Status:   StatusSuccess,
			Message:   "",
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
