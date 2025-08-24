package port

type Identifier interface {
	Generate(projectName string, organizationName string) string
}
