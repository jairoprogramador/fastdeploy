package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"fastdeploy/internal/domain/versioning/vos"
)

// CommandExecutor es una función helper para ejecutar comandos de forma sencilla.
func runGitCommand(ctx context.Context, repoPath string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error al ejecutar git %v: %s, %w", args, string(output), err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GitRepo es la implementación de la interfaz GitRepository.
type GitRepo struct{}

// NewGitRepo crea una nueva instancia de GitRepo.
func NewGitRepo() *GitRepo {
	return &GitRepo{}
}

// GetLastCommit obtiene el último commit de la rama actual.
func (r *GitRepo) GetLastCommit(ctx context.Context, repoPath string) (*vos.Commit, error) {
	// Formato: hash<|>subject<|>author_name<|>author_date_iso8601
	log, err := runGitCommand(ctx, repoPath, "log", "-1", "--pretty=format:%H<|>%s<|>%an<|>%ai")
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(log, "<|>", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("formato de log de git inesperado: %s", log)
	}

	date, err := time.Parse("2006-01-02 15:04:05 -0700", parts[3])
	if err != nil {
		return nil, fmt.Errorf("no se pudo parsear la fecha del commit: %w", err)
	}

	return &vos.Commit{
		Hash:    parts[0],
		Message: parts[1],
		Author:  parts[2],
		Date:    date,
	}, nil
}

// GetCommitsSinceTag obtiene la lista de commits desde un tag específico.
func (r *GitRepo) GetCommitsSinceTag(ctx context.Context, repoPath string, lastTag string) ([]*vos.Commit, error) {
	logRange := "HEAD"
	if lastTag != "" {
		logRange = fmt.Sprintf("%s..HEAD", lastTag)
	}

	// Formato: hash<|>subject
	logOutput, err := runGitCommand(ctx, repoPath, "log", "--pretty=format:%H<|>%s", logRange)
	if err != nil {
		// Si no hay commits, git log puede devolver un error. Lo tratamos como una lista vacía.
		if strings.Contains(err.Error(), "fatal: bad revision") {
			return []*vos.Commit{}, nil
		}
		return nil, err
	}

	if logOutput == "" {
		return []*vos.Commit{}, nil
	}

	lines := strings.Split(logOutput, "\n")
	commits := make([]*vos.Commit, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, "<|>", 2)
		if len(parts) != 2 {
			continue // Ignorar líneas mal formadas
		}
		commits = append(commits, &vos.Commit{
			Hash:    parts[0],
			Message: parts[1],
		})
	}
	return commits, nil
}

// GetLastSemverTag busca el último tag que coincida con un patrón semver.
func (r *GitRepo) GetLastSemverTag(ctx context.Context, repoPath string) (string, error) {
	// Obtenemos todos los tags ordenados por fecha de creación descendente
	tags, err := runGitCommand(ctx, repoPath, "tag", "-l", "v*.*.*", "--sort=-creatordate")
	if err != nil {
		return "", err
	}

	if tags == "" {
		// No se encontraron tags semver
		return "", nil
	}

	// El primer tag en la lista es el más reciente
	lastTag := strings.Split(tags, "\n")[0]
	return lastTag, nil
}
