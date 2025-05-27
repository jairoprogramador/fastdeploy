package entity

type StoreEntity struct {
	global map[string]string
	local  []map[string]string
}

func NewStoreEntity() *StoreEntity {
	return &StoreEntity{
		global: make(map[string]string),
		local:  make([]map[string]string, 0),
	}
}

func (s *StoreEntity) Initialize(variables []Variable) {
	s.global = make(map[string]string)

	for _, v := range variables {
		s.global[v.Name] = v.Value
	}
}

func (s *StoreEntity) AddVariableGlobal(name, value string) {
	if s.global == nil {
		s.global = make(map[string]string)
	}
	s.global[name] = value
}

func (s *StoreEntity) PushScope(variables []Variable) {
	scope := make(map[string]string)
	for _, v := range variables {
		scope[v.Name] = v.Value
	}
	s.local = append(s.local, scope)
}

func (s *StoreEntity) PopScope() {
	if len(s.local) > 0 {
		s.local = s.local[:len(s.local)-1]
	}
}

func (s *StoreEntity) Get(name string) string {
	for i := len(s.local) - 1; i >= 0; i-- {
		if value, exists := s.local[i][name]; exists {
			return value
		}
	}
	return s.global[name]
}
