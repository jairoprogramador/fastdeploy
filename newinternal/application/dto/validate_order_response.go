package dto

import "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"

type ValidateOrderResponse struct {
	FinalStep        string
	Environment      vos.Environment
}