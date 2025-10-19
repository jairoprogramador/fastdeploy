package dom

import (
	"crypto/sha256"
	"fmt"
	"sort"

	"github.com/jairoprogramador/fastdeploy/internal/domain/dom/services"
	"github.com/jairoprogramador/fastdeploy/internal/domain/dom/vos"
)

type ShaGenerator struct{}

func NewShaGenerator() services.ShaGenerator {
	return &ShaGenerator{}
}

func (g *ShaGenerator) GenerateProductID(name, organization string) vos.ProductID {
	data := fmt.Sprintf("name:%s|organization:%s", name, organization)
	hash := sha256.Sum256([]byte(data))
	return vos.ProductID(fmt.Sprintf("%x", hash))
}

func (g *ShaGenerator) GenerateProjectID(tech vos.Technology) vos.ProjectID {
	fields := []string{
		fmt.Sprintf("type:%s", tech.TypeTechnology()),
		fmt.Sprintf("solution:%s", tech.Solution()),
		fmt.Sprintf("stack:%s", tech.Stack()),
		fmt.Sprintf("infrastructure:%s", tech.Infrastructure()),
	}
	sort.Strings(fields)

	var combined string
	for _, field := range fields {
		combined += field + "|"
	}

	hash := sha256.Sum256([]byte(combined))
	return vos.ProjectID(fmt.Sprintf("%x", hash))
}
