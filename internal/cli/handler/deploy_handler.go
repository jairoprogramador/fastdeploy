package handler

type DeployHandler struct {}

func NewDeployHandler() *DeployHandler {
	return &DeployHandler{}
}

func (h *DeployHandler) Controller() error {
	return nil
}
