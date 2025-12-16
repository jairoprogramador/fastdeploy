package vos

type StepStatus string

const (
	Success StepStatus = "SUCCESS"
	Failure StepStatus = "FAILURE"
	Cached  StepStatus = "CACHED"
)

type ExecutionResult struct {
	Status     StepStatus
	Logs       string
	OutputVars VariableSet
	Error      error
}
