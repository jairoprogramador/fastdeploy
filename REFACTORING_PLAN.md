# Refactoring Plan for FastDeploy

This document outlines the plan for refactoring the FastDeploy project to better align with Domain-Driven Design (DDD) principles and Hexagonal Architecture (Ports & Adapters).

## Current Architecture Analysis

The current architecture already has some elements of DDD and Hexagonal Architecture:

1. **Domain Layer**: Contains entities, services, and repositories
2. **Infrastructure Layer**: Contains adapters and repository implementations
3. **Application Layer**: Contains use cases
4. **CLI Layer**: Contains commands and handlers

However, there are several areas for improvement:

1. The domain layer contains both core domain logic and ports (interfaces for external dependencies)
2. The bounded contexts are not clearly defined
3. Some domain logic is mixed with infrastructure concerns
4. The application layer is very thin and could be enhanced
5. The dependency injection is done manually in the main.go file

## Proposed Architecture

### 1. Domain Layer (Core)

The domain layer should contain only the core domain logic, including:

- **Entities**: Value objects that represent the core concepts of the domain
- **Aggregates**: Clusters of entities and value objects that are treated as a single unit
- **Domain Services**: Services that contain domain logic that doesn't naturally fit into entities
- **Domain Events**: Events that represent something that happened in the domain
- **Domain Exceptions**: Custom exceptions for domain-specific error cases

### 2. Application Layer

The application layer should coordinate the execution of use cases by:

- **Use Cases**: Orchestrating the flow of data to and from the domain layer
- **DTOs**: Data Transfer Objects for communication with external layers
- **Assemblers/Mappers**: Converting between domain entities and DTOs
- **Application Services**: Coordinating multiple use cases or complex flows

### 3. Infrastructure Layer

The infrastructure layer should provide implementations for the interfaces defined in the domain layer:

- **Repositories**: Implementations of domain repository interfaces
- **Adapters**: Implementations of domain port interfaces
- **External Services**: Integration with external systems
- **Persistence**: Database access and ORM configurations
- **Messaging**: Message queue implementations

### 4. Interface Layer

The interface layer should handle the interaction with the outside world:

- **CLI**: Command-line interface
- **API**: REST or GraphQL API
- **UI**: User interface components
- **Controllers/Handlers**: Handling requests and responses

## Bounded Contexts

Based on the analysis, we can identify the following bounded contexts:

1. **Project Management**: Handling project configuration and initialization
2. **Deployment**: Managing the deployment process
3. **Docker**: Handling Docker-related operations
4. **Configuration**: Managing application configuration

## Proposed Directory Structure

```
fastdeploy/
├── cmd/
│   └── fastdeploy/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── project/
│   │   │   ├── entity/
│   │   │   │   └── project.go
│   │   │   ├── repository/
│   │   │   │   └── project_repository.go
│   │   │   ├── service/
│   │   │   │   └── project_service.go
│   │   │   └── event/
│   │   │       └── project_events.go
│   │   ├── deployment/
│   │   │   ├── entity/
│   │   │   │   └── deployment.go
│   │   │   ├── repository/
│   │   │   │   └── deployment_repository.go
│   │   │   ├── service/
│   │   │   │   └── deployment_service.go
│   │   │   └── event/
│   │   │       └── deployment_events.go
│   │   ├── docker/
│   │   │   ├── entity/
│   │   │   │   └── container.go
│   │   │   ├── service/
│   │   │   │   └── docker_service.go
│   │   │   └── event/
│   │   │       └── docker_events.go
│   │   └── config/
│   │       ├── entity/
│   │       │   └── config.go
│   │       ├── repository/
│   │       │   └── config_repository.go
│   │       ├── service/
│   │       │   └── config_service.go
│   │       └── event/
│   │           └── config_events.go
│   ├── application/
│   │   ├── project/
│   │   │   ├── dto/
│   │   │   │   └── project_dto.go
│   │   │   ├── mapper/
│   │   │   │   └── project_mapper.go
│   │   │   └── usecase/
│   │   │       ├── initialize_project.go
│   │   │       └── start_project.go
│   │   ├── deployment/
│   │   │   ├── dto/
│   │   │   │   └── deployment_dto.go
│   │   │   ├── mapper/
│   │   │   │   └── deployment_mapper.go
│   │   │   └── usecase/
│   │   │       └── execute_deployment.go
│   │   └── docker/
│   │       ├── dto/
│   │       │   └── docker_dto.go
│   │       ├── mapper/
│   │       │   └── docker_mapper.go
│   │       └── usecase/
│   │           ├── start_container.go
│   │           └── check_container.go
│   ├── infrastructure/
│   │   ├── repository/
│   │   │   ├── yaml/
│   │   │   │   ├── yaml_project_repository.go
│   │   │   │   ├── yaml_deployment_repository.go
│   │   │   │   └── yaml_config_repository.go
│   │   │   └── file/
│   │   │       └── file_repository.go
│   │   ├── adapter/
│   │   │   ├── docker/
│   │   │   │   ├── local_docker_container.go
│   │   │   │   └── local_docker_image.go
│   │   │   ├── git/
│   │   │   │   └── local_git_request.go
│   │   │   ├── file/
│   │   │   │   └── os_file_controller.go
│   │   │   ├── path/
│   │   │   │   └── os_path_service.go
│   │   │   ├── command/
│   │   │   │   └── os_run_command.go
│   │   │   └── template/
│   │   │       └── text_docker_template.go
│   │   └── config/
│   │       └── app_config.go
│   └── interface/
│       ├── cli/
│       │   ├── command/
│       │   │   ├── root.go
│       │   │   ├── init.go
│       │   │   ├── start.go
│       │   │   └── deploy.go
│       │   ├── handler/
│       │   │   ├── init_handler.go
│       │   │   ├── start_handler.go
│       │   │   └── deploy_handler.go
│       │   └── presenter/
│       │       └── ui.go
│       └── di/
│           └── container.go
├── pkg/
│   ├── common/
│   │   ├── result/
│   │   │   ├── domain_result.go
│   │   │   └── infra_result.go
│   │   ├── logger/
│   │   │   ├── logger.go
│   │   │   ├── console_logger.go
│   │   │   └── file_logger.go
│   │   └── util/
│   │       └── id_generator.go
│   └── constant/
│       ├── default.go
│       ├── key.go
│       ├── messages.go
│       └── paths.go
└── test/
    ├── unit/
    │   └── domain/
    │       └── project/
    │           └── service/
    │               └── project_service_test.go
    └── integration/
        └── application/
            └── project/
                └── usecase/
                    └── initialize_project_test.go
```

## Implementation Plan

### Phase 1: Restructure the Domain Layer

1. Create bounded contexts for Project, Deployment, Docker, and Configuration
2. Move domain entities to their respective bounded contexts
3. Move domain services to their respective bounded contexts
4. Move repository interfaces to their respective bounded contexts
5. Create domain events for important state changes

### Phase 2: Enhance the Application Layer

1. Create DTOs for communication with external layers
2. Create mappers for converting between domain entities and DTOs
3. Implement use cases for each bounded context
4. Create application services for coordinating multiple use cases

### Phase 3: Refactor the Infrastructure Layer

1. Move repository implementations to the infrastructure layer
2. Move adapters to the infrastructure layer
3. Implement infrastructure services for external dependencies
4. Create a configuration service for managing application configuration

### Phase 4: Reorganize the Interface Layer

1. Move CLI commands to the interface layer
2. Move handlers to the interface layer
3. Create presenters for formatting output
4. Implement a dependency injection container

### Phase 5: Testing and Validation

1. Create unit tests for domain services
2. Create integration tests for use cases
3. Create end-to-end tests for CLI commands
4. Validate that the application still works as expected

## Benefits of the Refactoring

1. **Improved Maintainability**: Clear separation of concerns makes the code easier to understand and maintain
2. **Better Testability**: Domain logic is isolated and can be tested independently of infrastructure
3. **Enhanced Flexibility**: The application can be extended with new features without modifying existing code
4. **Reduced Coupling**: Dependencies are explicit and managed through interfaces
5. **Clearer Intent**: The code structure reflects the domain model, making it easier to understand the business logic

## Conclusion

This refactoring plan aims to improve the architecture of the FastDeploy project by applying DDD principles and Hexagonal Architecture. The proposed changes will make the code more maintainable, testable, and flexible, while preserving the existing functionality.