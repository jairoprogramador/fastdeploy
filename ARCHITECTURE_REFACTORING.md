# Refactorización de Arquitectura - FastDeploy

## Resumen de Cambios

Se ha realizado una refactorización importante de la arquitectura del proyecto para aplicar el **Principio de Inversión de Dependencias** y mejorar la separación de responsabilidades.

## Problemas Identificados

### Antes de la Refactorización

1. **Acoplamiento Directo al Sistema Operativo**: Los módulos del dominio (`config`, `project`, `initializer`) dependían directamente de librerías del sistema operativo como `os`, `os/user`, etc.

2. **Violación del Principio de Inversión de Dependencias**: El dominio (capa interna) dependía de detalles de implementación (capa externa).

3. **Dificultad para Testing**: Era imposible hacer unit tests sin depender del sistema de archivos real.

4. **Falta de Flexibilidad**: No se podía cambiar la implementación del sistema de archivos sin modificar el dominio.

## Solución Implementada

### 1. Creación de Interfaces en el Dominio

Se definieron interfaces claras en el dominio para las operaciones necesarias:

```go
// FileSystem interface para operaciones de archivos
type FileSystem interface {
    ReadFile(path string) ([]byte, error)
    WriteFile(path string, data []byte, perm uint32) error
    MkdirAll(path string, perm uint32) error
    IsNotExist(err error) bool
    Stat(name string) (os.FileInfo, error)
}

// UserSystem interface para operaciones de usuario
type UserSystem interface {
    Current() (*user.User, error)
}

// WorkingDirectory interface para operaciones del directorio de trabajo
type WorkingDirectory interface {
    Getwd() (string, error)
}
```

### 2. Implementaciones en Adaptadores

Se crearon implementaciones concretas en la capa de adaptadores:

```go
// OSFileSystem implementa FileSystem usando las librerías del sistema operativo
type OSFileSystem struct{}

// OSUserSystem implementa UserSystem usando las librerías del sistema operativo
type OSUserSystem struct{}

// OSWorkingDirectory implementa WorkingDirectory usando las librerías del sistema operativo
type OSWorkingDirectory struct{}
```

### 3. Servicios con Inyección de Dependencias

Se refactorizaron los servicios para recibir las dependencias por inyección:

```go
// ConfigService implementa la lógica de negocio de configuración
type ConfigService struct {
    fileSystem FileSystem
    userSystem UserSystem
}

// ProjectService implementa la lógica de negocio de proyecto
type ProjectService struct {
    fileSystem FileSystem
    workingDir WorkingDirectory
}

// ProjectInitializerService implementa la lógica de inicialización de proyectos
type ProjectInitializerService struct {
    projectService *ProjectService
    configService  *config.ConfigService
    fileSystem     FileSystem
}
```

### 4. Funciones Legacy para Compatibilidad

Se mantuvieron las funciones originales como wrappers para facilitar la migración gradual:

```go
// Funciones legacy para mantener compatibilidad (se eliminarán después de la migración)
func Save(configEntity ConfigEntity) error {
    service := NewConfigService(&filesystem.OSFileSystem{}, &filesystem.OSUserSystem{})
    return service.Save(configEntity)
}
```

## Estructura de Archivos Actualizada

```
internal/
├── adapters/
│   └── filesystem/
│       └── os_filesystem.go          # Implementaciones del sistema operativo
├── core/
│   └── domain/
│       ├── config/
│       │   ├── config.go             # Servicios con inyección de dependencias
│       │   └── config_entity.go      # Entidades del dominio
│       └── project/
│           ├── project.go            # Servicios con inyección de dependencias
│           ├── project_entity.go     # Entidades del dominio
│           └── initializer.go        # Servicios con inyección de dependencias
```

## Beneficios Obtenidos

### 1. **Separación de Responsabilidades**
- El dominio ya no conoce detalles de implementación del sistema operativo
- Las responsabilidades están claramente definidas y separadas

### 2. **Testabilidad Mejorada**
- Se pueden crear mocks de las interfaces para unit testing
- Los tests no dependen del sistema de archivos real

### 3. **Flexibilidad**
- Se pueden crear diferentes implementaciones (memoria, red, etc.)
- Fácil cambio de implementación sin modificar el dominio

### 4. **Principios SOLID Aplicados**
- **S**: Single Responsibility Principle - Cada servicio tiene una responsabilidad
- **O**: Open/Closed Principle - Abierto para extensión, cerrado para modificación
- **D**: Dependency Inversion Principle - El dominio no depende de detalles

## Ejemplo de Uso

### Antes (Acoplado)
```go
func Save(configEntity ConfigEntity) error {
    filePath, err := GetConfigFilePath()
    if err != nil {
        return err
    }
    
    data, err := yaml.Marshal(configEntity)
    if err != nil {
        return err
    }
    
    return os.WriteFile(filePath, data, 0644) // ❌ Dependencia directa
}
```

### Después (Desacoplado)
```go
func (cs *ConfigService) Save(configEntity ConfigEntity) error {
    filePath, err := cs.getConfigFilePath()
    if err != nil {
        return err
    }
    
    data, err := yaml.Marshal(configEntity)
    if err != nil {
        return err
    }
    
    return cs.fileSystem.WriteFile(filePath, data, 0644) // ✅ Dependencia inyectada
}
```

## Próximos Pasos

1. **Eliminar Funciones Legacy**: Una vez que todos los comandos usen la nueva arquitectura
2. **Agregar Tests Unitarios**: Crear mocks y tests para los servicios
3. **Implementar DI Container**: Considerar usar un contenedor de inyección de dependencias
4. **Documentar Interfaces**: Agregar documentación detallada a las interfaces

## Conclusión

Esta refactorización ha mejorado significativamente la arquitectura del proyecto, haciéndola más mantenible, testeable y flexible. El dominio ahora está correctamente aislado de los detalles de implementación, siguiendo los principios de Clean Architecture.
