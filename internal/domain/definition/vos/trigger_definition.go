package vos

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared"

type TriggerDefinition int

const (
	ScopeNone TriggerDefinition = iota
	ScopeCode
	ScopeRecipe
	ScopeVars
)

func (v TriggerDefinition) String() string {
	switch v {
	case ScopeCode:
		return shared.ScopeCode
	case ScopeRecipe:
		return shared.ScopeRecipe
	case ScopeVars:
		return shared.ScopeVars
	default:
		return shared.ScopeNone
	}
}

func TriggerFromString(s string) TriggerDefinition {
	switch s {
	case shared.ScopeCode:
		return ScopeCode
	case shared.ScopeRecipe:
		return ScopeRecipe
	case shared.ScopeVars:
		return ScopeVars
	default:
		return ScopeNone
	}
}
