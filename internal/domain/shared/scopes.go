package shared

const (
	ScopeCode = "code"
	ScopeRecipe = "recipe"
	ScopeVars = "vars"
	ScopeNone = "none"
)

func ScopesValid() []string {
	return []string{
		ScopeCode,
		ScopeRecipe,
		ScopeVars,
	}
}