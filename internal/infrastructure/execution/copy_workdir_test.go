package execution_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/execution"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyWorkdir_Success(t *testing.T) {
	// Arrange
	sourceDir := t.TempDir()
	destDir := t.TempDir()

	// Crear estructura de directorios y archivos de prueba
	require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "subdir", "file2.txt"), []byte("content2"), 0644))

	copier := execution.NewCopyWorkdir()

	// Act
	err := copier.Copy(context.Background(), sourceDir, destDir)

	// Assert
	require.NoError(t, err)

	// Verificar que los archivos y directorios existen en el destino
	assert.FileExists(t, filepath.Join(destDir, "file1.txt"))
	assert.DirExists(t, filepath.Join(destDir, "subdir"))
	assert.FileExists(t, filepath.Join(destDir, "subdir", "file2.txt"))

	// Verificar contenido de los archivos
	content1, err := os.ReadFile(filepath.Join(destDir, "file1.txt"))
	require.NoError(t, err)
	assert.Equal(t, "content1", string(content1))

	content2, err := os.ReadFile(filepath.Join(destDir, "subdir", "file2.txt"))
	require.NoError(t, err)
	assert.Equal(t, "content2", string(content2))
}

func TestCopyWorkdir_SourceNotExists(t *testing.T) {
	// Arrange
	sourceDir := filepath.Join(t.TempDir(), "nonexistent")
	destDir := t.TempDir()
	copier := execution.NewCopyWorkdir()

	// Act
	err := copier.Copy(context.Background(), sourceDir, destDir)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestCopyWorkdir_SourceIsNotDir(t *testing.T) {
	// Arrange
	sourceFile := filepath.Join(t.TempDir(), "file.txt")
	require.NoError(t, os.WriteFile(sourceFile, []byte("i am a file"), 0644))
	destDir := t.TempDir()
	copier := execution.NewCopyWorkdir()

	// Act
	err := copier.Copy(context.Background(), sourceFile, destDir)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no es un directorio")
}

func TestCopyWorkdir_ContextCancellation(t *testing.T) {
	// Arrange
	sourceDir := t.TempDir()
	destDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "file2.txt"), []byte("content2"), 0644))

	copier := execution.NewCopyWorkdir()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(2 * time.Nanosecond) // Asegurarse de que el contexto ya está cancelado

	// Act
	err := copier.Copy(ctx, sourceDir, destDir)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	// Verificamos que el directorio de destino esté vacío o parcialmente creado, pero no completo.
	entries, _ := os.ReadDir(destDir)
	assert.LessOrEqual(t, len(entries), 1, "El directorio de destino no debería estar completamente copiado")
}
