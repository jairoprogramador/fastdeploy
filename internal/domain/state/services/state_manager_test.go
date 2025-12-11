package services

import (
	"errors"
	"testing"
	"time"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
)

type mockStateRepository struct {
	GetFunc  func(workspacePath string, step vos.Step) (*aggregates.StateTable, error)
	SaveFunc func(workspacePath string, stateTable *aggregates.StateTable) error

	saveCalledWith *aggregates.StateTable
}

func (m *mockStateRepository) Get(workspacePath string, step vos.Step) (*aggregates.StateTable, error) {
	if m.GetFunc != nil {
		return m.GetFunc(workspacePath, step)
	}
	return aggregates.NewStateTable(step), nil
}

func (m *mockStateRepository) Save(workspacePath string, stateTable *aggregates.StateTable) error {
	m.saveCalledWith = stateTable
	if m.SaveFunc != nil {
		return m.SaveFunc(workspacePath, stateTable)
	}
	return nil
}

func newFingerprint(v string) vos.Fingerprint {
	fp, _ := vos.NewFingerprint(v)
	return fp
}

func newStep(v string) vos.Step {
	s, _ := vos.NewStep(v)
	return s
}

func newEnv(v string) vos.Environment {
	e, _ := vos.NewEnvironment(v)
	return e
}

func TestStateManager_HasStateChanged(t *testing.T) {
	step := newStep(vos.StepTest)
	env := newEnv("dev")

	fpCode1 := "code1"
	fpInst1 := "inst1"
	fpVars1 := "vars1"

	// Estado actual que usaremos para la comprobación
	currentState := vos.NewCurrentStateFingerprints(
		newFingerprint(fpCode1), newFingerprint(fpInst1), newFingerprint(fpVars1), env)
	// Una entrada histórica que coincide con el estado actual
	matchingEntry := aggregates.NewStateEntry(
		newFingerprint(fpCode1), newFingerprint(fpInst1), newFingerprint(fpVars1), env)

	// Una entrada histórica que coincide pero es demasiado antigua
	expiredMatchingEntry := aggregates.NewStateEntry(
		newFingerprint(fpCode1), newFingerprint(fpInst1), newFingerprint(fpVars1), env)
	// Forzamos su fecha de creación a ser de hace mucho tiempo para la prueba
	expiredMatchingEntry.SetCreatedAtForTesting(time.Now().Add(-31 * 24 * time.Hour))

	testCases := []struct {
		name         string
		repo         *mockStateRepository
		step         vos.Step
		currentState vos.CurrentStateFingerprints
		wantChanged  bool
		wantErr      bool
	}{
		{
			name: "debería devolver true (changed) si el repositorio devuelve un error (ej. no encontrado)",
			repo: &mockStateRepository{
				GetFunc: func(workspacePath string, step vos.Step) (*aggregates.StateTable, error) {
					return nil, errors.New("file not found")
				},
			},
			step:         step,
			currentState: currentState,
			wantChanged:  true,
			wantErr:      true,
		},
		{
			name: "debería devolver false (not changed) si se encuentra una coincidencia",
			repo: &mockStateRepository{
				GetFunc: func(workspacePath string, step vos.Step) (*aggregates.StateTable, error) {
					table := aggregates.NewStateTable(step)
					table.AddEntry(matchingEntry)
					return table, nil
				},
			},
			step:         step,
			currentState: currentState,
			wantChanged:  false,
			wantErr:      false,
		},
		{
			name: "debería devolver true (changed) si no se encuentra una coincidencia",
			repo: &mockStateRepository{
				GetFunc: func(workspacePath string, step vos.Step) (*aggregates.StateTable, error) {
					// Tabla con una entrada que no coincide
					table := aggregates.NewStateTable(step)
					nonMatchingEntry := aggregates.NewStateEntry(
						newFingerprint("other-code"),
						newFingerprint("inst1"),
						newFingerprint("vars1"), env)
					table.AddEntry(nonMatchingEntry)
					return table, nil
				},
			},
			step:         step,
			currentState: currentState,
			wantChanged:  true,
			wantErr:      false,
		},
		{
			name: "debería devolver true (changed) si la coincidencia ha expirado por tiempo",
			repo: &mockStateRepository{
				GetFunc: func(workspacePath string, step vos.Step) (*aggregates.StateTable, error) {
					table := aggregates.NewStateTable(step)
					table.AddEntry(expiredMatchingEntry)
					return table, nil
				},
			},
			step:         step,
			currentState: currentState,
			wantChanged:  true, // Debe ejecutar porque la entrada expiró
			wantErr:      false,
		},
		{
			name:         "debería devolver true (changed) y un error si el step es inválido",
			repo:         &mockStateRepository{},
			step:         newStep("invalid-step"),
			currentState: currentState,
			wantChanged:  true,
			wantErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// El fpService no se usa en estos métodos, así que podemos pasar nil
			sm := NewStateManager(tc.repo)

			// Usamos una política por defecto para las pruebas
			policy := vos.NewCachePolicy(0)
			gotChanged, err := sm.HasStateChanged("/fake/path", tc.step, tc.currentState, policy)

			if (err != nil) != tc.wantErr {
				t.Errorf("HasStateChanged() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if gotChanged != tc.wantChanged {
				t.Errorf("HasStateChanged() = %v, want %v", gotChanged, tc.wantChanged)
			}
		})
	}
}

func TestStateManager_UpdateState(t *testing.T) {
	step := newStep(vos.StepDeploy)
	env := newEnv("prod")

	currentState := vos.NewCurrentStateFingerprints(
		newFingerprint("code1"),
		newFingerprint("inst1"),
		newFingerprint("vars1"), env)

	testCases := []struct {
		name       string
		repo       *mockStateRepository
		expectErr  bool
		verifySave func(t *testing.T, repo *mockStateRepository)
	}{
		{
			name: "debería propagar el error si Get falla",
			repo: &mockStateRepository{
				GetFunc: func(workspacePath string, step vos.Step) (*aggregates.StateTable, error) {
					return nil, errors.New("I/O error")
				},
			},
			expectErr: true,
			verifySave: func(t *testing.T, repo *mockStateRepository) {
				if repo.saveCalledWith != nil {
					t.Error("Save no debería haber sido llamado si Get falla")
				}
			},
		},
		{
			name: "debería crear una nueva tabla si Get devuelve (nil, nil) y guardar la nueva entrada",
			repo: &mockStateRepository{
				GetFunc: func(workspacePath string, step vos.Step) (*aggregates.StateTable, error) {
					return nil, nil // Simula que el archivo no existe, pero sin error
				},
			},
			expectErr: false,
			verifySave: func(t *testing.T, repo *mockStateRepository) {
				if repo.saveCalledWith == nil {
					t.Fatal("Save no fue llamado")
				}
				if len(repo.saveCalledWith.Entries()) != 1 {
					t.Errorf("Se esperaba 1 entrada en la tabla guardada, pero se obtuvieron %d", len(repo.saveCalledWith.Entries()))
				}
				if repo.saveCalledWith.Step() != step {
					t.Errorf("El step de la tabla guardada es incorrecto")
				}
			},
		},
		{
			name: "debería añadir a una tabla existente y guardarla",
			repo: &mockStateRepository{
				GetFunc: func(workspacePath string, step vos.Step) (*aggregates.StateTable, error) {
					// Devuelve una tabla que ya tiene una entrada
					table := aggregates.NewStateTable(step)
					existingEntry := aggregates.NewStateEntry(
						newFingerprint("old"), newFingerprint("old"), newFingerprint("old"), env)
					table.AddEntry(existingEntry)
					return table, nil
				},
			},
			expectErr: false,
			verifySave: func(t *testing.T, repo *mockStateRepository) {
				if repo.saveCalledWith == nil {
					t.Fatal("Save no fue llamado")
				}
				if len(repo.saveCalledWith.Entries()) != 2 {
					t.Errorf("Se esperaban 2 entradas en la tabla guardada, pero se obtuvieron %d", len(repo.saveCalledWith.Entries()))
				}
			},
		},
		{
			name: "debería devolver un error si Save falla",
			repo: &mockStateRepository{
				SaveFunc: func(workspacePath string, stateTable *aggregates.StateTable) error {
					return errors.New("disk full")
				},
			},
			expectErr:  true,
			verifySave: nil, // No verificamos el guardado si se espera un error
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sm := NewStateManager(tc.repo)

			err := sm.UpdateState("/fake/path", step, currentState)

			if (err != nil) != tc.expectErr {
				t.Errorf("UpdateState() error = %v, wantErr %v", err, tc.expectErr)
			}

			if tc.verifySave != nil {
				tc.verifySave(t, tc.repo)
			}
		})
	}
}
