package execution

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/ports"
)

type CopyWorkdir struct{}

func NewCopyWorkdir() ports.CopyWorkdir {
	return &CopyWorkdir{}
}

func (c *CopyWorkdir) Copy(ctx context.Context, source, destination string) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("no se pudo obtener información de la fuente '%s': %w", source, err)
	}
	if !sourceInfo.IsDir() {
		return fmt.Errorf("la fuente '%s' no es un directorio", source)
	}

	return filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Comprobar si el contexto ha sido cancelado
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Calcular la ruta de destino
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return fmt.Errorf("no se pudo calcular la ruta relativa para '%s': %w", path, err)
		}
		destPath := filepath.Join(destination, relPath)

		fileInfo, err := d.Info()
		if err != nil {
			return fmt.Errorf("no se pudo obtener información de la entrada '%s': %w", path, err)
		}

		if d.IsDir() {
			// Crear el directorio en el destino con los mismos permisos que el origen
			if err := os.MkdirAll(destPath, fileInfo.Mode()); err != nil {
				return fmt.Errorf("no se pudo crear el directorio '%s': %w", destPath, err)
			}
		} else {
			// Copiar el archivo
			if err := copyFile(path, destPath, fileInfo.Mode()); err != nil {
				return fmt.Errorf("no se pudo copiar el archivo de '%s' a '%s': %w", path, destPath, err)
			}
		}

		return nil
	})
}

// copyFile copia el contenido de un archivo a otro de manera eficiente,
// estableciendo los permisos correctos en el momento de la creación.
func copyFile(src, dst string, perm fs.FileMode) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Usamos OpenFile para crear el archivo con los permisos correctos desde el principio.
	destFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
