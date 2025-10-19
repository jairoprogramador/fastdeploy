package vos

import "github.com/jairoprogramador/fastdeploy/internal/domain/shared"

type Trigger int

const (
	ScopeNone Trigger = iota
	ScopeCode
	ScopeRecipe
	ScopeVars
)

func (v Trigger) String() string {
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

func TriggerFromString(s string) Trigger{
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
