package ports

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"

type FingerprintService interface {
	FromFile(filePath string) (vos.Fingerprint, error)
	FromDirectory(dirPath string) (vos.Fingerprint, error)
}