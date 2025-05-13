package command

import (
	"deploy/internal/domain/engine"
	"deploy/internal/domain/model"
	"deploy/internal/infrastructure/filesystem"
	//"deploy/internal/infrastructure/repository"
	"deploy/internal/application/dto"
	//"deploy/internal/domain/service"

	"os"
	"log"
	"context"
	"time"
	"gopkg.in/yaml.v3"
)

func StartDeploy() *dto.ResponseDto {
	// Leer archivo YAML
	homeDir, err := filesystem.GetHomeDirectory()
	if err != nil {
		log.Fatalf("Error leyendo archivo: %v", err)
	}
	deploymentPath := filesystem.GetPath(homeDir, ".fastdeploy", "deployment.yaml")

    data, err := os.ReadFile(deploymentPath)
    if err != nil {
        log.Fatalf("Error leyendo archivo: %v", err)
    }

    // Parsear YAML
    var deployment model.Deployment
    if err := yaml.Unmarshal(data, &deployment); err != nil {
        log.Fatalf("Error parseando YAML: %v", err)
    }

    // Crear motor
    engine := engine.NewEngine()

    // Crear contexto con timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
    defer cancel()

    // Ejecutar despliegue
    if err := engine.Execute(ctx, &deployment); err != nil {
        log.Fatalf("%v", err)
    }

	return dto.GetDtoWithMessage("Despliegue completado exitosamente")
    //log.Println("Despliegue completado exitosamente")

	/* publishRepository := repository.GetPublishRepository()
	publishService := service.GetPublishService(publishRepository)

	responseBuild := publishService.Build()
	if responseBuild.Error != nil {
		return dto.GetNewResponseDtoFromModel(responseBuild) 
	}

	responsePackage := publishService.Package(responseBuild)
	if responsePackage.Error != nil {
		return dto.GetNewResponseDtoFromModel(responsePackage)
	}

	responseDeliver := publishService.Deliver(responsePackage)
	return dto.GetNewResponseDtoFromModel(responseDeliver) */
}
