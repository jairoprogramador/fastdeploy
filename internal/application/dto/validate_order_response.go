package dto

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/vos"

type ValidateOrderResponse struct {
	FinalStep        string
	Environment      vos.Environment
}