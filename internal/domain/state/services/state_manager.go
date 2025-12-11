package services

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
)

type StateManager struct {
	stateRepo ports.StateRepository
}

func NewStateManager(stateRepo ports.StateRepository) *StateManager {
	return &StateManager{
		stateRepo: stateRepo,
	}
}

func (sm *StateManager) HasStateChanged(
	workspacePath string,
	step vos.Step,
	currentState vos.CurrentStateFingerprints,
	policy vos.CachePolicy,
) (bool, error) {
	stateTable, err := sm.stateRepo.Get(workspacePath, step)
	if err != nil {
		return true, err
	}

	match, err := sm.findMatch(stateTable, currentState, policy)
	if err != nil {
		return true, err
	}

	return match == nil, nil
}

func (sm *StateManager) UpdateState(
	workspacePath string,
	step vos.Step,
	currentState vos.CurrentStateFingerprints,
) error {
	stateTable, err := sm.stateRepo.Get(workspacePath, step)
	if err != nil {
		return err
	}
	if stateTable == nil {
		stateTable = aggregates.NewStateTable(step)
	}

	newEntry := aggregates.NewStateEntry(
		currentState.Code(),
		currentState.Instruction(),
		currentState.Vars(),
		currentState.Environment(),
	)
	stateTable.AddEntry(newEntry)

	return sm.stateRepo.Save(workspacePath, stateTable)
}

func (sm *StateManager) findMatch(
	st *aggregates.StateTable,
	currentState vos.CurrentStateFingerprints,
	policy vos.CachePolicy,
) (*aggregates.StateEntry, error) {
	matcher, err := NewStateMatcherFactory(st.Step(), policy)
	if err != nil {
		return nil, err
	}

	for _, entry := range st.Entries() {
		if matcher.Match(entry, currentState) {
			return entry, nil
		}
	}
	return nil, nil
}
