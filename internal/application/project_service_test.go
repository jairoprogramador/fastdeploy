package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jairoprogramador/fastdeploy-core/internal/application"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/project/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/project/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeProjectRepository es un mock para el ProjectRepository.
type fakeProjectRepository struct {
	LoadFunc func(ctx context.Context, path string) (*ports.ProjectConfigDTO, error)
	SaveFunc func(ctx context.Context, path string, data *ports.ProjectConfigDTO) error

	saveCalled bool
}

func (f *fakeProjectRepository) Load(ctx context.Context, path string) (*ports.ProjectConfigDTO, error) {
	if f.LoadFunc != nil {
		return f.LoadFunc(ctx, path)
	}
	return nil, errors.New("LoadFunc no implementado")
}

func (f *fakeProjectRepository) Save(ctx context.Context, path string, data *ports.ProjectConfigDTO) error {
	f.saveCalled = true
	if f.SaveFunc != nil {
		return f.SaveFunc(ctx, path, data)
	}
	return nil
}

// fakeGitCloner es un mock para ClonerTemplate.
type fakeGitCloner struct {
	EnsureClonedFunc func(ctx context.Context, repoURL, ref, localPath string) error
	RunFunc          func(ctx context.Context, command string, workDir string) (*ports.CommandResultDTO, error)
}

func (f *fakeGitCloner) EnsureCloned(ctx context.Context, repoURL, ref, localPath string) error {
	if f.EnsureClonedFunc != nil {
		return f.EnsureClonedFunc(ctx, repoURL, ref, localPath)
	}
	return nil
}

func (f *fakeGitCloner) Run(ctx context.Context, command string, workDir string) (*ports.CommandResultDTO, error) {
	if f.RunFunc != nil {
		return f.RunFunc(ctx, command, workDir)
	}
	return &ports.CommandResultDTO{Output: "", ExitCode: 0}, nil
}

func newValidMockDTO(modifiers ...func(*ports.ProjectConfigDTO)) *ports.ProjectConfigDTO {
	// 1. Define los datos base y consistentes
	projectName := "test-project"
	projectOrganization := "fastdeploy"
	projectTeam := "shikigami"
	expectedID := vos.GenerateProjectID(projectName, projectOrganization, projectTeam)

	// 2. Crea el DTO base
	dto := &ports.ProjectConfigDTO{
		ID:           expectedID.String(),
		Name:         projectName,
		Organization: projectOrganization,
		Team:         projectTeam,
		Version:      "1.0.0",
		TemplateURL:  "https://github.com/jairo/template.git",
		TemplateRef:  "v1",
	}

	// 3. Aplica cualquier modificación específica del test
	for _, modifier := range modifiers {
		modifier(dto)
	}

	return dto
}

func TestProjectService_Initialize_Success(t *testing.T) {
	// --- Arrange ---
	ctx := context.Background()

	mockRepo := &fakeProjectRepository{
		LoadFunc: func(ctx context.Context, path string) (*ports.ProjectConfigDTO, error) {
			return newValidMockDTO(), nil // Usamos el helper sin modificaciones
		},
	}
	mockCloner := &fakeGitCloner{}
	service := application.NewProjectService(mockRepo, mockCloner)

	// --- Act ---
	project, err := service.Initialize(ctx, "/fake/path", "/fake/repos")

	// --- Assert ---
	require.NoError(t, err)
	require.NotNil(t, project)
	assert.Equal(t, "test-project", project.Data().Name())
	assert.Equal(t, "/fake/repos/template", project.TemplateLocalPath())
	assert.False(t, mockRepo.saveCalled, "Save no debería haber sido llamado porque el ID no cambió")
}

func TestProjectService_Initialize_IDChangesAndSaves(t *testing.T) {
	// --- Arrange ---
	ctx := context.Background()

	mockRepo := &fakeProjectRepository{
		LoadFunc: func(ctx context.Context, path string) (*ports.ProjectConfigDTO, error) {
			// Usamos el helper y le pasamos una función para modificar solo el ID
			return newValidMockDTO(func(dto *ports.ProjectConfigDTO) {
				dto.ID = "id-incorrecto"
			}), nil
		},
	}
	mockCloner := &fakeGitCloner{}
	service := application.NewProjectService(mockRepo, mockCloner)

	// --- Act ---
	project, err := service.Initialize(ctx, "/fake/path", "/fake/repos")

	// --- Assert ---
	require.NoError(t, err)
	require.NotNil(t, project)

	// Verificamos contra el ID que el helper habría calculado
	validDTO := newValidMockDTO()
	assert.Equal(t, validDTO.ID, project.ID().String(), "El ID del proyecto debería haberse actualizado")
	assert.True(t, mockRepo.saveCalled, "Save debería haber sido llamado porque el ID cambió")
}

func TestProjectService_Initialize_RepoLoadFails(t *testing.T) {
	// --- Arrange ---
	ctx := context.Background()
	expectedError := errors.New("failed to read file")

	mockRepo := &fakeProjectRepository{
		LoadFunc: func(ctx context.Context, path string) (*ports.ProjectConfigDTO, error) {
			return nil, expectedError
		},
	}
	mockCloner := &fakeGitCloner{}
	service := application.NewProjectService(mockRepo, mockCloner)

	// --- Act ---
	project, err := service.Initialize(ctx, "/fake/path", "/fake/repos")

	// --- Assert ---
	require.Error(t, err)
	assert.Nil(t, project)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestProjectService_Initialize_SaveFailsAfterIDChange(t *testing.T) {
	// --- Arrange ---
	ctx := context.Background()
	expectedError := errors.New("permission denied")

	mockRepo := &fakeProjectRepository{
		LoadFunc: func(ctx context.Context, path string) (*ports.ProjectConfigDTO, error) {
			return newValidMockDTO(func(dto *ports.ProjectConfigDTO) {
				dto.ID = "id-incorrecto"
			}), nil
		},
		SaveFunc: func(ctx context.Context, path string, data *ports.ProjectConfigDTO) error {
			return expectedError
		},
	}
	mockCloner := &fakeGitCloner{}
	service := application.NewProjectService(mockRepo, mockCloner)

	// --- Act ---
	project, err := service.Initialize(ctx, "/fake/path", "/fake/repos")

	// --- Assert ---
	require.Error(t, err)
	assert.Nil(t, project)
	assert.True(t, mockRepo.saveCalled)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestProjectService_Initialize_ClonerFails(t *testing.T) {
	// --- Arrange ---
	ctx := context.Background()
	expectedError := errors.New("git command not found")

	mockRepo := &fakeProjectRepository{
		LoadFunc: func(ctx context.Context, path string) (*ports.ProjectConfigDTO, error) {
			return newValidMockDTO(), nil // ID válido, la carga es exitosa
		},
	}
	mockCloner := &fakeGitCloner{
		EnsureClonedFunc: func(ctx context.Context, repoURL, ref, localPath string) error {
			return expectedError
		},
	}
	service := application.NewProjectService(mockRepo, mockCloner)

	// --- Act ---
	project, err := service.Initialize(ctx, "/fake/path", "/fake/repos")

	// --- Assert ---
	require.Error(t, err)
	assert.Nil(t, project)
	assert.Contains(t, err.Error(), expectedError.Error())
}
