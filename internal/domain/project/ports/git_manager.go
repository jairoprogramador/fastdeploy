package ports

type GitManager interface {
	Clone(url string, nameRepository string) error
	IsCloned(nameRepository string) (bool, error)
}
