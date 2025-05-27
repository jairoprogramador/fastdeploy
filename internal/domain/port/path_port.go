package port

type PathPort interface {
	GetFullPathDockerComposeTemplate() string

	GetFullPathDockerCompose() string

	GetFullPathDockerfileTemplate() string

	GetFullPathDockerfile() string

	GetPathDockerDirectory() string

	GetPathProjectFile() string

	GetFullPathConfigFile() string

	GetFullPathDeploymentFile() string

	GetFullPathLoggerFile() string

	GetProjectName() (string, error)

	GetHomeDirectory() string

	GetRelativePathFromHome(absolutePath string) string
}
