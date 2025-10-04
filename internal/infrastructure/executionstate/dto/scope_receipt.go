package dto

import "time"

type FingerprintDTO struct {
	Value string `yaml:"value"`
}

type ScopeReceiptDTO struct {
	CodeFingerprint        FingerprintDTO
	EnvironmentFingerprint FingerprintDTO
	CreatedAt              time.Time
}

type ScopeReceiptHistoryDTO struct {
	Receipts []*ScopeReceiptDTO
}
