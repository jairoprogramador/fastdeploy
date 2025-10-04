package dto

import "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"

type ValidateOrderResponse struct {
	FinalStep        string
	Environment      vos.Environment
}