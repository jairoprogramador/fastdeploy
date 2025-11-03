package dto

type ExecutorRequest struct {
	environment      string
	finalStepName    string
	pathProject      string
	skippedStepNames map[string]struct{}
}

func NewExecutorRequest(
	environment, finalStepName string,
	pathProject string,
	skippedStepNames map[string]struct{}) ExecutorRequest {
	return ExecutorRequest{
		environment:      environment,
		finalStepName:    finalStepName,
		pathProject:      pathProject,
		skippedStepNames: skippedStepNames,
	}
}

func (r *ExecutorRequest) Environment() string {
	return r.environment
}

func (r *ExecutorRequest) FinalStepName() string {
	return r.finalStepName
}

func (r *ExecutorRequest) SkippedStepNames() map[string]struct{} {
	return r.skippedStepNames
}

func (r *ExecutorRequest) PathProject() string {
	return r.pathProject
}
