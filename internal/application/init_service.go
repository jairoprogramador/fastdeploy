package application

import (
	"fmt"
	"time"

	"github.com/jairoprogramador/fastdeploy/internal/application/dto"
	"github.com/jairoprogramador/fastdeploy/internal/application/ports"
	domaggregates "github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
	domports "github.com/jairoprogramador/fastdeploy/internal/domain/dom/ports"
	domservices "github.com/jairoprogramador/fastdeploy/internal/domain/dom/services"
	domvos "github.com/jairoprogramador/fastdeploy/internal/domain/dom/vos"
)

// InitService es el servicio de aplicación que coordina la creación
// del Modelo de Objeto de Despliegue (dom.yaml).
type InitService struct {
	userInput ports.UserInputProvider
	idGen     domservices.IDGenerator
	domRepo   domports.DOMRepository
}

// NewInitService crea una nueva instancia de InitService.
func NewInitService(
	userInput ports.UserInputProvider,
	idGen domservices.IDGenerator,
	domRepo domports.DOMRepository) *InitService {
	return &InitService{
		userInput: userInput,
		idGen:     idGen,
		domRepo:   domRepo,
	}
}

// InitializeDOM orquesta la creación interactiva (o por defecto) del dom.yaml.
func (s *InitService) InitializeDOM(req dto.InitRequest) (*domaggregates.DeploymentObjectModel, error) {
	defaultProductName := req.WorkingDirectory
	defaultProjectName := req.WorkingDirectory
	defaultProjectVersion := time.Now().Format("2006.01.02")

	var err error
	var val string

	// Product
	productOrg := "fastdeploy"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Nombre de la compañia", productOrg)
		if err != nil { return nil, err }
		productOrg = val
	}

	productTeam := "zenin"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Nombre del equipo a cargo del producto", productTeam)
		if err != nil { return nil, err }
		productTeam = val
	}
	productName := defaultProductName
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Nombre del producto", productName)
		if err != nil { return nil, err }
		productName = val
	}

	projectName := defaultProjectName
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Nombre del proyecto", projectName)
		if err != nil { return nil, err }
		projectName = val
	}
	projectVersion := defaultProjectVersion
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Versión del proyecto", projectVersion)
		if err != nil { return nil, err }
		projectVersion = val
	}

	projectTeam := "shikigami"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Nombre del equipo a cargo del proyecto", projectTeam)
		if err != nil { return nil, err }
		projectTeam = val
	}

	templateUrl := "https://github.com/jairoprogramador/mydeploy.git"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "URL de la plantilla", templateUrl)
		if err != nil { return nil, err }
		templateUrl = val
	}

	templateRef := "main"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Version de la plantilla", templateRef)
		if err != nil { return nil, err }
		templateRef = val
	}

	techType := "backend"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Tipo de proyecto", techType)
		if err != nil { return nil, err }
		techType = val
	}
	techSolution := "microservice"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Solución tecnológica", techSolution)
		if err != nil { return nil, err }
		techSolution = val
	}
	techStack := "springboot"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Stack tecnológico", techStack)
		if err != nil { return nil, err }
		techStack = val
	}
	techInfrastructure := "azure"
	if !req.SkipPrompt {
		val, err = s.userInput.Prompt(req.Ctx, "Infraestructura tecnológica", techInfrastructure)
		if err != nil { return nil, err }
		techInfrastructure = val
	}

	// Generar IDs basados en los datos recopilados
	productID := s.idGen.GenerateProductID(productName, productOrg)

	techVO, err := domvos.NewTechnology(techType, techSolution, techStack, techInfrastructure)
	if err != nil { return nil, fmt.Errorf("error al crear el VO de tecnología: %w", err) }

	projectID := s.idGen.GenerateProjectID(*techVO)

	// Crear los VOs
	productVO, err := domvos.NewProduct(productID, productName, "my product", productTeam, productOrg)
	if err != nil { return nil, fmt.Errorf("error al crear el VO de producto: %w", err) }

	projectVO, err := domvos.NewProject(projectID, projectName, projectVersion, "my project", projectTeam)
	if err != nil { return nil, fmt.Errorf("error al crear el VO de proyecto: %w", err) }

	templateVO, err := domvos.NewTemplate(templateUrl, templateRef)
	if err != nil { return nil, fmt.Errorf("error al crear el VO de plantilla: %w", err) }

	// Crear el Agregado Raíz
	dom := domaggregates.NewDeploymentObjectModel(productVO, projectVO, templateVO, techVO)

	// --- Persistencia ---
	if err := s.domRepo.Save(req.Ctx, dom); err != nil {
		return nil, fmt.Errorf("error al guardar el archivo .fastdeploy/dom.yaml: %w", err)
	}

	fmt.Println("✅ Archivo .fastdeploy/dom.yaml creado exitosamente.")
	return dom, nil
}
