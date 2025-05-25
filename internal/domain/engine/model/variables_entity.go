package model

type Variables struct {
	Global []Variable `yaml:"global"`
	Local  []Variable `yaml:"local"`
}
