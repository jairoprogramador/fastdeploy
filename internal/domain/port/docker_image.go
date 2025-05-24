package port

type DockerImage interface {
	CreateDockerfile() error
}
