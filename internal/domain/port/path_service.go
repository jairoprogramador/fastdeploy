// Package port defines interfaces for external dependencies of the domain layer.
package port

// PathService defines the interface for path-related operations.
type PathService interface {
	// GetFullPathDockerComposeTemplate returns the full path to the Docker Compose template file.
	GetFullPathDockerComposeTemplate() string

	// GetFullPathDockerCompose returns the full path to the Docker Compose file.
	GetFullPathDockerCompose() string

	// GetFullPathDockerfileTemplate returns the full path to the Dockerfile template.
	GetFullPathDockerfileTemplate() string

	// GetFullPathDockerfile returns the full path to the Dockerfile.
	GetFullPathDockerfile() string

	// GetPathDockerDirectory returns the path to the Docker directory.
	GetPathDockerDirectory() string

	// GetPathProjectFile returns the path to the project file.
	GetPathProjectFile() string

	// GetFullPathConfigFile returns the full path to the configuration file.
	GetFullPathConfigFile() string

	// GetFullPathDeploymentFile returns the full path to the deployment file.
	GetFullPathDeploymentFile() string

	// GetFullPathLoggerFile returns the full path to the logger file.
	GetFullPathLoggerFile() string

	// GetProjectName returns the name of the current project.
	GetProjectName() (string, error)

	// GetHomeDirectory returns the home directory of the current user.
	GetHomeDirectory() string

	// GetRelativePathFromHome returns the relative path from the home directory.
	GetRelativePathFromHome(absolutePath string) string
}
