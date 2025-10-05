package mapper

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
)

// VerificationsToDomain convierte un slice de strings a un slice de VOs VerificationType.
func VerificationsToDomain(dtoStrings []string) ([]vos.VerificationType, error) {
	verifications := make([]vos.VerificationType, 0, len(dtoStrings))
	for _, s := range dtoStrings {
		verificationType, err := vos.VerificationTypeFromString(s)
		if err != nil {
			return nil, fmt.Errorf("error al convertir el tipo de verificación: %w", err)
		}
		verifications = append(verifications, verificationType)
	}
	return verifications, nil
}
