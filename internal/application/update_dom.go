package application

import (
	"context"
	"fmt"

	appPor "github.com/jairoprogramador/fastdeploy/internal/application/ports"

	domAgg "github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
	domPor "github.com/jairoprogramador/fastdeploy/internal/domain/dom/ports"
	domSer "github.com/jairoprogramador/fastdeploy/internal/domain/dom/services"

	staPor "github.com/jairoprogramador/fastdeploy/internal/domain/state/ports"
	staVos "github.com/jairoprogramador/fastdeploy/internal/domain/state/vos"
)

type UpdateDOMService struct {
	domRepository   domPor.DomRepository
	stateRepository staPor.FingerprintRepository
	idGenerator     domSer.ShaGenerator
	userInput       appPor.UserInputProvider
}

func NewUpdateDOMService(
	domRepository domPor.DomRepository,
	stateRepository staPor.FingerprintRepository,
	idGenerator domSer.ShaGenerator,
	userInput appPor.UserInputProvider) *UpdateDOMService {
	return &UpdateDOMService{
		domRepository:   domRepository,
		stateRepository: stateRepository,
		idGenerator:     idGenerator,
		userInput:       userInput,
	}
}

func (s *UpdateDOMService) Update(
	ctx context.Context, domModel *domAgg.DeploymentObjectModel) error {

	if domModel.IsModified(s.idGenerator) {

		err := domModel.UpdateIDs(s.idGenerator)
		if err != nil {
			return err
		}
		executionState, err := s.stateRepository.FindStep("supply")
		if err != nil {
			return err
		}
		if executionState != nil && executionState.ExistsFingerprint(staVos.ScopeVars) {
			fmt.Println("⚠️ Se han detectado cambios en .fastdeploy/dom.yaml que afectan a la identidad del proyecto.")
			confirmed, err := s.userInput.Confirm(ctx, "¿Continuar? Esto podría causar cambios en la infraestructura existente.", false)
			if err != nil || !confirmed {
				return fmt.Errorf("operación cancelada")
			}
		}
		if err := s.domRepository.Save(domModel); err != nil {
			return err
		}
		fmt.Println("✅ IDs del proyecto actualizados en .fastdeploy/dom.yaml.")
	}
	return nil
}
