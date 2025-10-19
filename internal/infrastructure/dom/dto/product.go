package dto

type ProductDTO struct {
	ID           string `yaml:"product_id"`
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	Team         string `yaml:"team"`
	Organization string `yaml:"organization"`
}
