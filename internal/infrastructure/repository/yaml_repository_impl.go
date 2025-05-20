package repository

import (
	"deploy/internal/domain/repository"

	"gopkg.in/yaml.v3"
)

type YamlRepositoryImpl struct {
	fileRepository repository.FileRepository
}

func NewYamlRepositoryImpl(fileRepo repository.FileRepository) repository.YamlRepository {
	return &YamlRepositoryImpl{
		fileRepository: fileRepo,
	}
}

func (st *YamlRepositoryImpl) Load(pathFile string, out any) error {
	file, err := st.fileRepository.OpenFile(pathFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(out)

	return err
}

func (st *YamlRepositoryImpl) Save(pathFile string, data any) error {
	file, err := st.fileRepository.CreateFile(pathFile)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	return encoder.Encode(data)
}
