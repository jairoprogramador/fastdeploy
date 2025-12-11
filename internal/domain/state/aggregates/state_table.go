package aggregates

import (
	"sort"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
)

const maxEntries = 5

type StateTable struct {
	step    vos.Step
	entries []*StateEntry
}

func NewStateTable(step vos.Step) *StateTable {
	return &StateTable{
		step:    step,
		entries: make([]*StateEntry, 0, maxEntries),
	}
}

func LoadStateTable(step vos.Step, entries []*StateEntry) *StateTable {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreatedAt().Before(entries[j].CreatedAt())
	})

	if len(entries) > maxEntries {
		entries = entries[len(entries)-maxEntries:]
	}

	return &StateTable{
		step:    step,
		entries: entries,
	}
}

func (st *StateTable) Entries() []*StateEntry {
	return st.entries
}

func (st *StateTable) Step() vos.Step {
	return st.step
}

func (st *StateTable) AddEntry(newEntry *StateEntry) {
	idx := sort.Search(len(st.entries), func(i int) bool {
		return st.entries[i].CreatedAt().After(newEntry.CreatedAt())
	})

	st.entries = append(st.entries[:idx], append([]*StateEntry{newEntry}, st.entries[idx:]...)...)

	if len(st.entries) > maxEntries {
		st.entries = st.entries[1:]
	}
}
