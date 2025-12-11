package matchers

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
)

// BaseMatcher encapsula la lógica de comparación común para todos los matchers.
type BaseMatcher struct{}

// matchCommon comprueba los fingerprints que son comunes a todas las estrategias.
func (b *BaseMatcher) matchCommon(entry *aggregates.StateEntry, current vos.CurrentStateFingerprints) bool {
	return entry.Instruction().Equals(current.Instruction()) &&
		entry.Vars().Equals(current.Vars())
}
