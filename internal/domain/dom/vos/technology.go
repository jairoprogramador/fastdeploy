package vos

const (
	DefaultStack = "springboot"
	DefaultInfrastructure = "azure"
)

type Technology struct {
	stack          string
	infrastructure string
}

func NewTechnology(stack, infrastructure string) Technology {
	if stack == "" {
		stack = DefaultStack
	}
	if infrastructure == "" {
		infrastructure = DefaultInfrastructure
	}
	return Technology{
		stack:          stack,
		infrastructure: infrastructure,
	}
}

func (t Technology) Stack() string { return t.stack }
func (t Technology) Infrastructure() string { return t.infrastructure }