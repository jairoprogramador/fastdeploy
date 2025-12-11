package ports

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
)

type StateMatcher interface {
	Match(entry *aggregates.StateEntry, current vos.CurrentStateFingerprints) bool
}
