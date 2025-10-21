package application

import (
	"fmt"

	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"

	domAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/aggregates"
	domPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/ports"
	domSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/services"
	domVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/vos"

	shaVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared/vos"
)

type InitService struct {
	userInput appPor.UserInputProvider
	idGen     domSer.ShaGenerator
	domRepo   domPor.DomRepository
}

func NewInitService(
	userInput appPor.UserInputProvider,
	idGen domSer.ShaGenerator,
	domRepo domPor.DomRepository) *InitService {
	return &InitService{
		userInput: userInput,
		idGen:     idGen,
		domRepo:   domRepo,
	}
}

func (s *InitService) Run(req appDto.InitRequest) (*domAgg.DeploymentObjectModel, error) {
	defaultProductName := req.WorkingDirectory
	defaultProjectName := req.WorkingDirectory
	defaultProjectVersion := "1.0.0"

	var err error
	var val string

	// Product
	productOrg := "fastdeploy"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Nombre de la compañia", productOrg)
		if err != nil {
			return nil, err
		}
		productOrg = val
	}

	productTeam := "zenin"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Nombre del equipo a cargo del producto", productTeam)
		if err != nil {
			return nil, err
		}
		productTeam = val
	}
	productName := defaultProductName
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Nombre del producto", productName)
		if err != nil {
			return nil, err
		}
		productName = val
	}

	projectName := defaultProjectName
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Nombre del proyecto", projectName)
		if err != nil {
			return nil, err
		}
		projectName = val
	}
	projectVersion := defaultProjectVersion
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Versión del proyecto", projectVersion)
		if err != nil {
			return nil, err
		}
		projectVersion = val
	}

	projectTeam := "shikigami"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Nombre del equipo a cargo del proyecto", projectTeam)
		if err != nil {
			return nil, err
		}
		projectTeam = val
	}

	templateUrl := "https://github.com/jairoprogramador/mydeploy.git"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "URL de la plantilla", templateUrl)
		if err != nil {
			return nil, err
		}
		templateUrl = val
	}

	templateRef := "main"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Version de la plantilla", templateRef)
		if err != nil {
			return nil, err
		}
		templateRef = val
	}

	techType := "backend"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Tipo de proyecto", techType)
		if err != nil {
			return nil, err
		}
		techType = val
	}
	techSolution := "microservice"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Solución tecnológica", techSolution)
		if err != nil {
			return nil, err
		}
		techSolution = val
	}
	techStack := "springboot"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Stack tecnológico", techStack)
		if err != nil {
			return nil, err
		}
		techStack = val
	}
	techInfrastructure := "azure"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Infraestructura tecnológica", techInfrastructure)
		if err != nil {
			return nil, err
		}
		techInfrastructure = val
	}

	productID := s.idGen.GenerateProductID(productName, productOrg)

	techVO, err := domVos.NewTechnology(techType, techSolution, techStack, techInfrastructure)
	if err != nil {
		return nil, fmt.Errorf("error al crear el VO de tecnología: %w", err)
	}

	projectID := s.idGen.GenerateProjectID(techVO)

	// Crear los VOs
	productVO, err := domVos.NewProduct(productID, productName, "my product", productTeam, productOrg)
	if err != nil {
		return nil, fmt.Errorf("error al crear el VO de producto: %w", err)
	}

	projectVO, err := domVos.NewProject(projectID, projectName, projectVersion, "my project", projectTeam)
	if err != nil {
		return nil, fmt.Errorf("error al crear el VO de proyecto: %w", err)
	}

	templateVO, err := shaVos.NewTemplateSource(templateUrl, templateRef)
	if err != nil {
		return nil, fmt.Errorf("error al crear el VO de plantilla: %w", err)
	}

	// Crear el Agregado Raíz
	dom := domAgg.NewDeploymentObjectModel(productVO, projectVO, templateVO, techVO)

	// --- Persistencia ---
	if err := s.domRepo.Save(dom); err != nil {
		return nil, fmt.Errorf("error al guardar el archivo .fastdeploy/dom.yaml: %w", err)
	}

	fmt.Println("✅ Archivo .fastdeploy/dom.yaml creado exitosamente.")
	return dom, nil
}
