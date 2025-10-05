package dto

// DOMDTO es la representación directa del archivo dom.yaml.
type DOMDTO struct {
	Product    ProductDTO    `yaml:"product"`
	Project    ProjectDTO    `yaml:"project"`
	Template   TemplateDTO   `yaml:"template"`
	Technology TechnologyDTO `yaml:"technology"`
}

type ProductDTO struct {
	ID           string `yaml:"product_id"`
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	Team         string `yaml:"team"`
	Organization string `yaml:"organization"`
}

type ProjectDTO struct {
	ID          string `yaml:"project_id"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Team        string `yaml:"team"`
}

type TemplateDTO struct {
	RepositoryURL string `yaml:"repository_url"`
	Ref           string `yaml:"ref"`
}

type TechnologyDTO struct {
	Type           string `yaml:"type"`
	Solution       string `yaml:"solution"`
	Stack          string `yaml:"stack"`
	Infrastructure string `yaml:"infrastructure"`
}
