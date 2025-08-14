# Servicios y Repositorios - FastDeploy

Esta documentación describe la nueva estructura de servicios y repositorios implementada en FastDeploy.

## Estructura de Capas

### 1. Servicios de Dominio (Domain Services)

Ubicación: `internal/domain/services/`

#### ConfigService
- **Función**: Orquesta la lógica de configuración
- **Métodos**:
  - `Load(repo ConfigRepository) -> Config`: Carga configuración desde archivo
  - `Save(repo ConfigRepository, config Config)`: Guarda configuración en archivo
  - `createDefault() -> Config`: Crea configuración por defecto

#### ProjectService
- **Función**: Orquesta la lógica de proyectos
- **Métodos**:
  - `Initialize(configService, projectRepo, projectName) -> Project`: Inicializa nuevo proyecto
  - `Load(repo ProjectRepository) -> Project`: Carga proyecto desde archivo
  - `saveProject(repo ProjectRepository, project Project)`: Guarda proyecto en archivo

### 2. Repositorios (Infrastructure/Repositories)

Ubicación: `internal/infrastructure/repositories/`

#### ConfigRepository
- **Constante**: `CONFIG_FILE_NAME = "config.yaml"`
- **Métodos**:
  - `Load() -> dict`: Lee datos del archivo de configuración
  - `Save(data dict)`: Escribe datos en archivo de configuración

#### ProjectRepository
- **Constante**: `PROJECT_FILE_NAME = "deploy.yaml"`
- **Métodos**:
  - `Load() -> dict`: Lee datos del archivo de proyecto
  - `Save(data dict)`: Escribe datos en archivo de proyecto

### 3. Servicios de Aplicación (Application Services)

Ubicación: `internal/application/`

#### ProjectApplicationService
- **Función**: Coordina operaciones de alto nivel
- **Métodos**:
  - `CreateProject(projectName string)`: Inicializa nuevo proyecto
  - `GetProject() -> Project`: Obtiene datos de proyecto existente

## Flujo de Trabajo

### Crear un Nuevo Proyecto

1. **ProjectApplicationService.CreateProject()**
   - Instancia repositorios y servicios de dominio
   - Llama a ProjectService.Initialize()

2. **ProjectService.Initialize()**
   - Usa ConfigService para obtener configuración
   - Crea entidades del proyecto (ProjectID, ProjectName, etc.)
   - Guarda el proyecto usando ProjectRepository

3. **ConfigService.Load()**
   - Intenta cargar configuración desde archivo
   - Si no existe, crea configuración por defecto

### Obtener un Proyecto Existente

1. **ProjectApplicationService.GetProject()**
   - Instancia repositorio y servicio de dominio
   - Llama a ProjectService.Load()

2. **ProjectService.Load()**
   - Usa ProjectRepository para cargar datos
   - Convierte datos a entidad Project

## Archivos de Configuración

- **config.yaml**: Contiene configuración global del proyecto
- **deploy.yaml**: Contiene datos específicos del proyecto

## Uso de Ejemplo

```go
// Crear servicio de aplicación
appService := application.NewProjectApplicationService()

// Crear nuevo proyecto
err := appService.CreateProject("mi-proyecto")
if err != nil {
    log.Printf("Error: %v", err)
}

// Obtener proyecto existente
project, err := appService.GetProject()
if err != nil {
    log.Printf("Error: %v", err)
}
```

## Notas de Implementación

- Los repositorios manejan la persistencia en archivos YAML
- Los servicios de dominio contienen la lógica de negocio
- Los servicios de aplicación coordinan las operaciones
- Se manejan errores de I/O apropiadamente
- Las entidades de dominio se crean con valores por defecto cuando es necesario
