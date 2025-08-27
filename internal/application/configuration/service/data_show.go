package service

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entity"
)

type DataShow struct {
	fileReader Reader
}

func NewDataShow(fileReader Reader) *DataShow {
	return &DataShow{
		fileReader: fileReader,
	}
}

func (ds *DataShow) Show() (entity.Configuration, error) {
	config, err := ds.fileReader.Read()
	if err != nil {
		return entity.Configuration{}, err
	}

	fmt.Println("Configuración de FastDeploy")
	fmt.Printf("\tOrganización: %s\n", config.GetNameOrganization().Value())
	fmt.Printf("\tEquipo: %s\n", config.GetTeam().Value())
	fmt.Printf("\tTecnología: %s\n", config.GetTechnology().Value())
	fmt.Printf("\tRepositorio.url: %s\n", config.GetRepository().GetURL().Value())

	return config, nil
}