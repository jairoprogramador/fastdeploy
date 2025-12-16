package services

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"sort"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
)

// ShaFingerprintGenerator es una implementación de FingerprintGenerator que usa SHA256.
type ShaFingerprintGenerator struct{}

// NewShaFingerprintGenerator crea una nueva instancia del generador de fingerprints.
func NewShaFingerprintGenerator() ports.FingerprintService {
	return &ShaFingerprintGenerator{}
}

func (s *ShaFingerprintGenerator) FromString(data string) (vos.Fingerprint, error) {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return vos.NewFingerprint(hash)
}

func (s *ShaFingerprintGenerator) FromFile(filePath string) (vos.Fingerprint, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return vos.Fingerprint{}, err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return vos.Fingerprint{}, err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	return vos.NewFingerprint(hash)
}

func (s *ShaFingerprintGenerator) FromMap(data map[string]string) (vos.Fingerprint, error) {
	if len(data) == 0 {
		return s.FromString("")
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	hasher := sha256.New()
	for _, k := range keys {
		hasher.Write([]byte(k))
		hasher.Write([]byte(data[k]))
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	return vos.NewFingerprint(hash)
}
