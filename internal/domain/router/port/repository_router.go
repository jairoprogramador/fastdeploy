package port

type RepositoryRouter interface {
	BuildRoute(paths ...string) string
}
