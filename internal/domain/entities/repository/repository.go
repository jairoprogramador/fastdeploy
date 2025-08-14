package repository

type Repository struct {
	url RepositoryURL
}

func NewRepository(url RepositoryURL) Repository {
	return Repository{
		url: url,
	}
}

func (r Repository) GetURL() RepositoryURL {
	return r.url
}

func (r Repository) GetName() (RepositoryName, error) {
	return ExtractFromURL(r.url)
}

func (r Repository) IsValid() bool {
	return r.url.Value() != ""
}

func (r Repository) Equals(other Repository) bool {
	return r.url.Value() == other.url.Value()
}
