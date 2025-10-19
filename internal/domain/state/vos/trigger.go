package vos

type Trigger int

const (
	ScopeNone Trigger = iota
	ScopeCode
	ScopeRecipe
	ScopeVars
)

func NewTrigger(value int) Trigger{
	switch value {
	case int(ScopeCode):
		return ScopeCode
	case int(ScopeRecipe):
		return ScopeRecipe
	case int(ScopeVars):
		return ScopeVars
	}
	return ScopeNone
}

func (v Trigger) Int() int {
	switch v {
	case ScopeCode:
		return int(ScopeCode)
	case ScopeRecipe:
		return int(ScopeRecipe)
	case ScopeVars:
		return int(ScopeVars)
	default:
		return int(ScopeNone)
	}
}
