package config

type ConfigValidator interface {
	Validate(config ConfigEntity) error
}
