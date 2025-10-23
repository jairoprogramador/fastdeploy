package dom

import (
	"crypto/sha256"
	"fmt"
	"sort"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/services"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/vos"
)

type ShaGeneratorID struct{}

func NewShaGeneratorID() services.GeneratorID {
	return &ShaGeneratorID{}
}

func (g *ShaGeneratorID) ProjectID(config *aggregates.Config) vos.ProjectID {
	fields := []string{
		fmt.Sprintf("template:%s", config.Template().NameTemplate()),
		fmt.Sprintf("stack:%s", config.Technology().Stack()),
		fmt.Sprintf("infrastructure:%s", config.Technology().Infrastructure()),
	}
	sort.Strings(fields)

	var combined string
	for _, field := range fields {
		combined += field + "|"
	}

	hash := sha256.Sum256([]byte(combined))
	return vos.ProjectID(fmt.Sprintf("%x", hash))
}
