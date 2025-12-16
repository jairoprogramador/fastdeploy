package ports

// FileSystem defines the port for file system operations, abstracting the
// underlying implementation details from the domain logic.
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
}
