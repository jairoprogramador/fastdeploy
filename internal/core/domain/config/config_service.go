package config

type ConfigService interface {
	Save(configEntity ConfigEntity) error
	Load() (*ConfigEntity, error)
	Exists() bool
	Delete() error
}
