package config

type ConfigRepository interface {
	Save(configEntity ConfigEntity) error
	Load() (*ConfigEntity, error)
	Exists() bool
	Delete() error
}
