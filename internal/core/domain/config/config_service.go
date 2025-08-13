package config

type ServiceReader interface {
	Load() (*ConfigEntity, error)
	Exists() bool
}

type ServiceWriter interface {
	Save(configEntity ConfigEntity) error
	Delete() error
}

type ConfigService interface {
	ServiceReader
	ServiceWriter
}
