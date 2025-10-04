package aggregates

import (
	"sort"

	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/vos"
)

const maxHistorySize = 3

type ScopeReceiptHistory struct {
	receipts []*ScopeReceipt
}

func NewScopeReceiptHistory() (*ScopeReceiptHistory, error) {
	return &ScopeReceiptHistory{
		receipts: make([]*ScopeReceipt, 0),
	}, nil
}

func RehydrateScopeReceiptHistory(receipts []*ScopeReceipt) *ScopeReceiptHistory {
	return &ScopeReceiptHistory{
		receipts: receipts,
	}
}

func (h *ScopeReceiptHistory) AddReceipt(receipt *ScopeReceipt) {
	h.receipts = append(h.receipts, receipt)

	sort.Slice(h.receipts, func(i, j int) bool {
		return h.receipts[i].CreatedAt().After(h.receipts[j].CreatedAt())
	})

	if len(h.receipts) > maxHistorySize {
		h.receipts = h.receipts[:maxHistorySize]
	}
}

func (h *ScopeReceiptHistory) findMatch(codeFp, envFp vos.Fingerprint) *ScopeReceipt {
	for _, receipt := range h.receipts {
		codeMatch := (receipt.CodeFingerprint() == codeFp)
		envMatch := (receipt.EnvironmentFingerprint() == envFp)

		if codeMatch && envMatch {
			return receipt
		}
	}
	return nil
}

func (h *ScopeReceiptHistory) FindMatchCode(codeFingerprint vos.Fingerprint) *ScopeReceipt {
	return h.findMatch(codeFingerprint, vos.Fingerprint{})
}

func (h *ScopeReceiptHistory) FindMatchEnvironment(environmentFingerprint vos.Fingerprint) *ScopeReceipt {
	return h.findMatch(vos.Fingerprint{}, environmentFingerprint)
}

func (h *ScopeReceiptHistory) Receipts() []*ScopeReceipt {
	receiptsCopy := make([]*ScopeReceipt, len(h.receipts))
	copy(receiptsCopy, h.receipts)
	return receiptsCopy
}
