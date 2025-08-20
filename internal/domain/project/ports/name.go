package ports

type Name interface {
	GetName() (string, error)
}
