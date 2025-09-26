package port

type ContextPort interface {
	Load(pathFileContext string) (map[string]string, error)
	Save(pathFileContext string, context map[string]string) error
}
