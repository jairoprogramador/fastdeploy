package git

type GitService interface {
	Clone(repositoryURL string) error
	IsCloned(repositoryURL string) bool
}
