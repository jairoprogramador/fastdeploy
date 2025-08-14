package application

import (
	"fmt"
	"log"
)

// ExampleUsage muestra cómo usar los servicios creados
func ExampleUsage() {
	// Crear el servicio de aplicación
	appService := NewProjectApplicationService()

	// Ejemplo 1: Crear un nuevo proyecto
	fmt.Println("=== Creando nuevo proyecto ===")
	err := appService.CreateProject("mi-proyecto")
	if err != nil {
		log.Printf("Error creando proyecto: %v", err)
	} else {
		fmt.Println("Proyecto creado exitosamente")
	}

	// Ejemplo 2: Obtener un proyecto existente
	fmt.Println("\n=== Obteniendo proyecto existente ===")
	project, err := appService.GetProject()
	if err != nil {
		log.Printf("Error obteniendo proyecto: %v", err)
	} else {
		fmt.Printf("Proyecto obtenido: %s\n", project.GetFullName())
	}
}
