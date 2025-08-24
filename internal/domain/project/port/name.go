package port

type Name interface {
	GetName() (string, error)
}
