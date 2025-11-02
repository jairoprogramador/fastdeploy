package dto

type RunParams struct {
	environment string
	stepName    string
}

func NewRunParams(environment, stepName string) RunParams {
	return RunParams{
		environment: environment,
		stepName:    stepName,
	}
}

func (r *RunParams) Environment() string {
	return r.environment
}

func (r *RunParams) StepName() string {
	return r.stepName
}
