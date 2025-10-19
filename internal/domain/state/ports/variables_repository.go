package ports

type VariablesRepository interface {
	FindByStepName(stepName string) (map[string]string, error)
	Save(stepName string, vars map[string]string) error
}