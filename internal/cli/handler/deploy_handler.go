package handler

import "github.com/jairoprogramador/fastdeploy/internal/domain/model"

type DeployHandler struct{}

func NewDeployHandler() *DeployHandler {
	return &DeployHandler{}
}

func (h *DeployHandler) Controller() model.DomainResultEntity {
	return model.NewResultApp("controller deploy not implemented")
}
