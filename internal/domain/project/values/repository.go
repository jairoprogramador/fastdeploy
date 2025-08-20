package values

type Repository struct {
	url UrlRepository
}

func NewRepository(url UrlRepository) Repository {
	return Repository{
		url: url,
	}
}

func NewDefaultRepository() Repository {
	defaultRepo := NewDefaultUrlRepository()
	return NewRepository(defaultRepo)
}

func (r Repository) GetURL() UrlRepository {
	return r.url
}

func (r Repository) GetName() (NameRepository, error) {
	return NewNameRepository(r.url.ExtractNameRepository())
}

func (r Repository) Equals(other Repository) bool {
	return r.url.Equals(other.url)
}
