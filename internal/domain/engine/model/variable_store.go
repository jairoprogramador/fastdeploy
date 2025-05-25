package model

type VariableStore struct {
	global map[string]string
	local  []map[string]string
}

func NewVariableStore() *VariableStore {
	return &VariableStore{
		global: make(map[string]string),
		local:  make([]map[string]string, 0),
	}
}

func (s *VariableStore) Initialize(variables []Variable) {
	s.global = make(map[string]string)

	for _, v := range variables {
		s.global[v.Name] = v.Value
	}
}

func (s *VariableStore) AddVariableGlobal(name, value string) {
	if s.global == nil {
		s.global = make(map[string]string)
	}
	s.global[name] = value
}

func (s *VariableStore) PushScope(variables []Variable) {
	scope := make(map[string]string)
	for _, v := range variables {
		scope[v.Name] = v.Value
	}
	s.local = append(s.local, scope)
}

func (s *VariableStore) PopScope() {
	if len(s.local) > 0 {
		s.local = s.local[:len(s.local)-1]
	}
}

func (s *VariableStore) Get(name string) string {
	for i := len(s.local) - 1; i >= 0; i-- {
		if value, exists := s.local[i][name]; exists {
			return value
		}
	}
	return s.global[name]
}
