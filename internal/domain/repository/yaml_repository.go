package repository

type YamlRepository interface {
	Load(pathFile string, out any) error
	Save(pathFile string, data any) error
} 
