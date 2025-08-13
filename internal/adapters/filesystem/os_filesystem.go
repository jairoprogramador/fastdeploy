package filesystem

import (
	"os"
	"os/user"
)

// FileSystem interface para operaciones de archivos
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm uint32) error
	MkdirAll(path string, perm uint32) error
	IsNotExist(err error) bool
	Stat(name string) (os.FileInfo, error)
}

// UserSystem interface para operaciones de usuario
type UserSystem interface {
	Current() (*user.User, error)
}

// WorkingDirectory interface para operaciones del directorio de trabajo
type WorkingDirectory interface {
	Getwd() (string, error)
}

// OSFileSystem implementa FileSystem usando las librerías del sistema operativo
type OSFileSystem struct{}

// ReadFile lee un archivo del sistema de archivos
func (fs *OSFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile escribe datos a un archivo
func (fs *OSFileSystem) WriteFile(path string, data []byte, perm uint32) error {
	return os.WriteFile(path, data, os.FileMode(perm))
}

// MkdirAll crea directorios recursivamente
func (fs *OSFileSystem) MkdirAll(path string, perm uint32) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

// IsNotExist verifica si un error indica que el archivo no existe
func (fs *OSFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// Stat obtiene información sobre un archivo
func (fs *OSFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// OSUserSystem implementa UserSystem usando las librerías del sistema operativo
type OSUserSystem struct{}

// Current obtiene el usuario actual
func (us *OSUserSystem) Current() (*user.User, error) {
	return user.Current()
}

// OSWorkingDirectory implementa WorkingDirectory usando las librerías del sistema operativo
type OSWorkingDirectory struct{}

// Getwd obtiene el directorio de trabajo actual
func (wd *OSWorkingDirectory) Getwd() (string, error) {
	return os.Getwd()
}
