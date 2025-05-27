package condition

type TypeCondition string

const (
	NotEmpty TypeCondition = "not_empty"
	Empty    TypeCondition = "empty"
	Equals   TypeCondition = "equals"
	Contains TypeCondition = "contains"
	Matches  TypeCondition = "matches"
)
