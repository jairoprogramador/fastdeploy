package dto

type OutputDto struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	//Format      string `yaml:"format,omitempty"`
	Regex       string `yaml:"regex,omitempty"`
	//MatchKey    string `yaml:"matchKey,omitempty"`
}