package vos

const (
	DefaultStateBackend = "local"
	DefaultStateURL     = ""
)

type State struct {
	backend string
	url     string
}

func NewState(backend, url string) State {
	if backend == "" {
		backend = DefaultStateBackend
	}
	if url == "" {
		url = DefaultStateURL
	}
	return State{backend: backend, url: url}
}

func (s State) Backend() string { return s.backend }
func (s State) URL() string { return s.url }