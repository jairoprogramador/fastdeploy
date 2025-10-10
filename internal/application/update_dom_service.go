package application

import (
	"context"
	"fmt"

	appports "github.com/jairoprogramador/fastdeploy/internal/application/ports"
	"github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
	domports "github.com/jairoprogramador/fastdeploy/internal/domain/dom/ports"
	"github.com/jairoprogramador/fastdeploy/internal/domain/dom/services"
	executionstateports "github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/ports"
)

type UpdateDOMService struct {
	domRepository     domports.DOMRepository
	scopeRepository executionstateports.ScopeRepository
	idGenerator       services.IDGenerator
	userInput         appports.UserInputProvider
}

func NewUpdateDOMService(
	domRepository domports.DOMRepository,
	scopeRepository executionstateports.ScopeRepository,
	idGenerator services.IDGenerator,
	userInput appports.UserInputProvider) *UpdateDOMService {
	return &UpdateDOMService{
		domRepository:     domRepository,
		scopeRepository: scopeRepository,
		idGenerator:       idGenerator,
		userInput:         userInput,
	}
}

func (s *UpdateDOMService) Update(
	ctx context.Context, domModel *aggregates.DeploymentObjectModel) error {
	isModified, err := domModel.VerifyAndUpdateIDs(s.idGenerator)
	if err != nil {
		return err
	}

	if isModified {
		history, _ := s.scopeRepository.FindStepStateHistory("supply")
		if history != nil && len(history.Receipts()) > 0 {
			fmt.Println("⚠️ Se han detectado cambios en .fastdeploy/dom.yaml que afectan a la identidad del proyecto.")
			confirmed, err := s.userInput.Confirm(ctx, "¿Continuar? Esto podría causar cambios en la infraestructura existente.", false)
			if err != nil || !confirmed {
				return fmt.Errorf("operación cancelada")
			}
		}
		if err := s.domRepository.Save(ctx, domModel); err != nil {
			return err
		}
		fmt.Println("✅ IDs del proyecto actualizados en .fastdeploy/dom.yaml.")
	}
	return nil
}
