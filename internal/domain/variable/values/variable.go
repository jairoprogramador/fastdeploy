package values

type VariableValue struct {
	name  string
	value string
}

func NewVariable(name, value string) VariableValue {
	return VariableValue{name: name, value: value}
}

func (v VariableValue) GetName() string {
	return v.name
}

func (v VariableValue) GetValue() string {
	return v.value
}
