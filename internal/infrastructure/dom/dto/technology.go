package dto

type TechnologyDTO struct {
	Type           string `yaml:"type"`
	Solution       string `yaml:"solution"`
	Stack          string `yaml:"stack"`
	Infrastructure string `yaml:"infrastructure"`
}
