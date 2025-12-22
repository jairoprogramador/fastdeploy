package ports

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"

type FingerprintService interface {
	//FromString(data string) (vos.Fingerprint, error)
	FromFile(filePath string) (vos.Fingerprint, error)
	FromDirectory(dirPath string) (vos.Fingerprint, error)
	//FromMap(data map[string]string) (vos.Fingerprint, error)
}