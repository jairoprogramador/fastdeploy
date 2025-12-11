package matchers

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
)

type PackageStateMatcher struct {
	BaseMatcher
}

func (m *PackageStateMatcher) Match(entry *aggregates.StateEntry, current vos.CurrentStateFingerprints) bool {
	if !m.matchCommon(entry, current) {
		return false
	}
	// environment is ignored for package step
	return entry.Code().Equals(current.Code())
}
