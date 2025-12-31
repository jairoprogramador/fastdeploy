package aggregates

import (
	"sort"
)

const maxEntries = 5

type StateTable struct {
	name    string
	entries []*StateEntry
}

func NewStateTable(name string) *StateTable {
	return &StateTable{
		name:    name,
		entries: make([]*StateEntry, 0, maxEntries),
	}
}

func LoadStateTable(name string, entries []*StateEntry) *StateTable {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreatedAt().Before(entries[j].CreatedAt())
	})

	if len(entries) > maxEntries {
		entries = entries[len(entries)-maxEntries:]
	}

	return &StateTable{
		name:    name,
		entries: entries,
	}
}

func (st *StateTable) Entries() []*StateEntry {
	return st.entries
}

func (st *StateTable) Name() string {
	return st.name
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
