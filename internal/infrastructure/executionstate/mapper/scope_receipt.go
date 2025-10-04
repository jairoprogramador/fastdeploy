package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/aggregates"
	executionstatevos "github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/vos"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executionstate/dto"
)

func ScopeReceiptHistoryToDTO(history *aggregates.ScopeReceiptHistory) *dto.ScopeReceiptHistoryDTO {
	receiptDTOs := make([]*dto.ScopeReceiptDTO, 0, len(history.Receipts()))
	for _, r := range history.Receipts() {
		receiptDTOs = append(receiptDTOs, &dto.ScopeReceiptDTO{
			CodeFingerprint:        dto.FingerprintDTO{Value: r.CodeFingerprint().String()},
			EnvironmentFingerprint: dto.FingerprintDTO{Value: r.EnvironmentFingerprint().String()},
			CreatedAt:              r.CreatedAt(),
		})
	}
	return &dto.ScopeReceiptHistoryDTO{
		Receipts: receiptDTOs,
	}
}

func ScopeReceiptHistoryToDomain(dto dto.ScopeReceiptHistoryDTO) *aggregates.ScopeReceiptHistory {
	receipts := make([]*aggregates.ScopeReceipt, 0, len(dto.Receipts))
	for _, rDTO := range dto.Receipts {
		codeFingerprint, _ := executionstatevos.NewFingerprint(rDTO.CodeFingerprint.Value)
		environmentFingerprint, _ := executionstatevos.NewFingerprint(rDTO.EnvironmentFingerprint.Value)

		receipt, _ := aggregates.RehydrateScopeReceipt(
			codeFingerprint,
			environmentFingerprint,
			rDTO.CreatedAt,
		)
		receipts = append(receipts, receipt)
	}
	return aggregates.RehydrateScopeReceiptHistory(receipts)
}
