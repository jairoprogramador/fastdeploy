package versioning

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoGitRepository(t *testing.T) {
	// Usamos un repo en disco temporal para poder pasar la ruta.
	// La implementación actual usa PlainOpen, que requiere una ruta.
	tmpDir := t.TempDir()
	repo, err := git.PlainInit(tmpDir, false)
	require.NoError(t, err)

	w, err := repo.Worktree()
	require.NoError(t, err)

	author := &object.Signature{Name: "Test", Email: "test@test.com", When: time.Now()}

	// Helper para crear un archivo y hacer commit
	commitFile := func(msg, content string) plumbing.Hash {
		t.Helper()
		err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(content), 0644)
		require.NoError(t, err)
		_, err = w.Add("test.txt")
		require.NoError(t, err)
		hash, err := w.Commit(msg, &git.CommitOptions{Author: author})
		require.NoError(t, err)
		return hash
	}

	// Commit 1 + tag v0.1.0
	c1 := commitFile("fix: bug", "content v1")
	_, err = repo.CreateTag("v0.1.0", c1, nil)
	require.NoError(t, err)

	// Commit 2 (sin tag)
	commitFile("feat: new feature", "content v2")

	// Commit 3 + tag v1.0.0 y otro no semántico
	c3 := commitFile("feat: another one", "content v3")
	_, err = repo.CreateTag("v1.0.0", c3, nil)
	require.NoError(t, err)
	_, err = repo.CreateTag("beta", c3, nil)
	require.NoError(t, err)

	// Commit 4 (HEAD)
	lastCommitMsg := "docs: update readme"
	c4 := commitFile(lastCommitMsg, "content v4")

	// --- Start Tests ---
	repoService := NewGoGitRepository()
	ctx := context.Background()

	t.Run("GetLastCommit", func(t *testing.T) {
		lastCommit, err := repoService.GetLastCommit(ctx, tmpDir)
		require.NoError(t, err)
		assert.Equal(t, c4.String(), lastCommit.Hash)
		assert.Contains(t, lastCommit.Message, lastCommitMsg)
	})

	t.Run("GetLastSemverTag", func(t *testing.T) {
		lastTag, err := repoService.GetLastSemverTag(ctx, tmpDir)
		require.NoError(t, err)
		assert.Equal(t, "v1.0.0", lastTag)
	})

	t.Run("GetCommitsSinceTag", func(t *testing.T) {
		t.Run("desde un tag específico", func(t *testing.T) {
			// Debería devolver 1 commit (el c4)
			commits, err := repoService.GetCommitsSinceTag(ctx, tmpDir, "v1.0.0")
			require.NoError(t, err)
			// La iteración devuelve el más nuevo primero, pero el tag está en c3,
			// por lo que solo c4 debería aparecer.
			assert.Len(t, commits, 1, "Debería haber 1 commit desde el tag v1.0.0")
			assert.Equal(t, c4.String(), commits[0].Hash)
		})

		t.Run("sin tag previo", func(t *testing.T) {
			// Debería devolver todos los commits
			commits, err := repoService.GetCommitsSinceTag(ctx, tmpDir, "")
			require.NoError(t, err)
			assert.Len(t, commits, 4, "Debería devolver todos los commits")
		})
	})
}
