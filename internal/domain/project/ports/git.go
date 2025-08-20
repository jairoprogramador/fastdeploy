package ports

type Git interface {
	Clone(url string, nameRepository string) error
	IsCloned(nameRepository string) (bool, error)
}