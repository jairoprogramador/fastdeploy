package port

type ProjectRouter interface {
	BuildRoute(paths ...string) string
}
