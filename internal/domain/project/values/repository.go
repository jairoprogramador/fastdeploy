package values

type Repository struct {
	url UrlRepository
	version VersionRepository
}

func NewRepository(url UrlRepository, version VersionRepository) Repository {
	return Repository{
		url: url,
		version: version,
	}
}

func NewDefaultRepository() Repository {
	defaultRepo := NewDefaultUrlRepository()
	defaultVersion := NewDefaultVersionRepository()
	return NewRepository(defaultRepo, defaultVersion)
}

func (r Repository) GetURL() UrlRepository {
	return r.url
}

func (r Repository) GetVersion() VersionRepository {
	return r.version
}

func (r Repository) GetName() (NameRepository, error) {
	return NewNameRepository(r.url.ExtractNameRepository())
}

func (r Repository) Equals(other Repository) bool {
	return r.url.Equals(other.url)
}
