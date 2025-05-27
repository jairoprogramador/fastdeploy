package handler

import "github.com/jairoprogramador/fastdeploy/pkg/common/result"

type DeployHandler struct{}

func NewDeployHandler() *DeployHandler {
	return &DeployHandler{}
}

func (h *DeployHandler) Controller() result.DomainResult {
	return result.NewResultApp("controller deploy not implemented")
}
