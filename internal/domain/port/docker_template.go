package port

type DockerTemplate interface {
	GetContent(pathTemplate string, params any) (string, error)
}
