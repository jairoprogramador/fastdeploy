package config

type RepositoryReader interface {
	Load() (*ConfigEntity, error)
	Exists() bool
}

type RepositoryWriter interface {
	Save(configEntity ConfigEntity) error
	Delete() error
}

type ConfigRepository interface {
	RepositoryReader
	RepositoryWriter
}
