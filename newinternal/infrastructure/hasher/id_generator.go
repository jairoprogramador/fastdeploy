package hasher

import (
	"crypto/sha256"
	"fmt"
	"sort"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/dom/services"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/dom/vos"
)

// IDGenerator implementa la interfaz services.IDGenerator.
// Utiliza SHA-256 para crear hashes deterministas.
type IDGenerator struct{}

// NewIDGenerator crea una nueva instancia del generador de IDs.
func NewIDGenerator() services.IDGenerator {
	return &IDGenerator{}
}

// GenerateProductID crea un hash a partir del nombre y la organización.
func (g *IDGenerator) GenerateProductID(name, organization string) vos.ProductID {
	data := fmt.Sprintf("name:%s|organization:%s", name, organization)
	hash := sha256.Sum256([]byte(data))
	return vos.ProductID(fmt.Sprintf("%x", hash))
}

// GenerateProjectID crea un hash a partir de los campos de la tecnología.
// Ordenamos los campos para asegurar que el hash sea estable sin importar el orden.
func (g *IDGenerator) GenerateProjectID(tech vos.Technology) vos.ProjectID {
	fields := []string{
		fmt.Sprintf("type:%s", tech.TypeTechnology()),
		fmt.Sprintf("solution:%s", tech.Solution()),
		fmt.Sprintf("stack:%s", tech.Stack()),
		fmt.Sprintf("infrastructure:%s", tech.Infrastructure()),
	}
	sort.Strings(fields) // CRÍTICO: asegura un hash determinista.

	var combined string
	for _, field := range fields {
		combined += field + "|"
	}

	hash := sha256.Sum256([]byte(combined))
	return vos.ProjectID(fmt.Sprintf("%x", hash))
}
