package vos

type VariableSet map[string]string

func (vs VariableSet) Clone() VariableSet {
	clone := make(VariableSet, len(vs))
	for k, v := range vs {
		clone[k] = v
	}
	return clone
}

func (vs VariableSet) Add(key, value string) {
	vs[key] = value
}

func (vs VariableSet) AddAll(other VariableSet) {
	for k, v := range other {
		vs[k] = v
	}
}