package entity

type Variables struct {
	Global []Variable `yaml:"global"`
	Local  []Variable `yaml:"local"`
}
