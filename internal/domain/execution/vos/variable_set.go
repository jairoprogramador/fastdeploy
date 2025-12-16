package vos

type VariableSet map[string]string

func (vs VariableSet) Clone() VariableSet {
	clone := make(VariableSet, len(vs))
	for k, v := range vs {
		clone[k] = v
	}
	return clone
}
