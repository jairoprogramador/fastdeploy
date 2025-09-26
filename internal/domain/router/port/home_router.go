package port

type HomeRouter interface {
	GetEnvironmentVariable(nameVariable string) string
	GetCurrentUserDir() (string, error)
	BuildRoute(paths ...string) string
	GetPathWorkdir() (string, error)
}
