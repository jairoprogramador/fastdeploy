package port

type WorkdirPort interface {
	Copy(sourcePath string, destinationPath string) error
}
