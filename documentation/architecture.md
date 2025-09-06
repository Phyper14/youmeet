# Architecture Guide

## Overview

YouMeet follows the **Hexagonal Architecture** (also known as Ports and Adapters pattern) to achieve separation of concerns, testability, and maintainability.

## Architecture Layers

### Domain Layer (`internal/domain/`)

The core business logic and entities. This layer is independent of external concerns.

**Components:**
- `entities.go` - Domain entities (User, Service, Appointment)

**Responsibilities:**
- Define business entities and their properties
- Contain business rules and invariants
- Independent of external frameworks and databases

### Ports Layer (`internal/ports/`)

Defines interfaces that the application uses to interact with external systems.

**Components:**
- `services.go` - Application service interfaces
- `repositories.go` - Repository interfaces for data persistence

**Responsibilities:**
- Define contracts for external dependencies
- Provide abstraction for application services
- Enable dependency inversion

### Application Layer (`internal/application/`)

Contains application-specific business logic and orchestrates domain entities.

**Components:**
- `appointment_service.go` - Appointment business logic implementation

**Responsibilities:**
- Implement application use cases
- Orchestrate domain entities
- Depend on ports, not concrete implementations

### Adapters Layer (`internal/adapters/`)

Implements the ports interfaces and handles external concerns.

**Components:**
- `http_handler.go` - HTTP REST API adapter
- `memory_repository.go` - In-memory data storage adapter

**Responsibilities:**
- Implement port interfaces
- Handle external system integration
- Convert between external formats and domain models

### Infrastructure Layer (`cmd/`)

Application entry point and dependency injection.

**Components:**
- `main.go` - Application bootstrap and dependency wiring

**Responsibilities:**
- Initialize and wire dependencies
- Start the application
- Configure external systems

## Dependency Flow

```
┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│  HTTP Handler   │
└─────────────────┘    └─────────────────┘
                                │
                                ▼
┌─────────────────┐    ┌─────────────────┐
│   Repository    │◀───│ Application     │
│   (Memory)      │    │ Service         │
└─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   Domain        │
                       │   Entities      │
                       └─────────────────┘
```

## Key Principles

### Dependency Inversion

- High-level modules don't depend on low-level modules
- Both depend on abstractions (interfaces)
- Application layer depends on ports, not concrete adapters

### Separation of Concerns

- **Domain**: Business logic and entities
- **Application**: Use cases and workflows  
- **Adapters**: External system integration
- **Infrastructure**: Configuration and bootstrapping

### Testability

- Each layer can be tested in isolation
- Dependencies can be easily mocked through interfaces
- Business logic is independent of external systems

## Package Structure

```
internal/
├── domain/           # Business entities and rules
│   └── entities.go
├── ports/            # Interface definitions
│   ├── services.go   # Application service interfaces
│   └── repositories.go # Repository interfaces
├── application/      # Use case implementations
│   └── appointment_service.go
└── adapters/         # External system adapters
    ├── http_handler.go    # HTTP REST API
    └── memory_repository.go # In-memory storage
```

## Benefits

1. **Maintainability**: Clear separation makes code easier to understand and modify
2. **Testability**: Each component can be tested independently
3. **Flexibility**: Easy to swap implementations (e.g., database, web framework)
4. **Independence**: Business logic is isolated from external concerns
5. **Scalability**: Architecture supports growth and complexity

## Design Patterns Used

- **Ports and Adapters**: Main architectural pattern
- **Dependency Injection**: Used in main.go for wiring components
- **Repository Pattern**: For data access abstraction
- **Service Layer**: For application logic organization

## Future Considerations

- Add database adapter (PostgreSQL, MySQL)
- Implement authentication and authorization
- Add caching layer
- Implement event-driven architecture
- Add API versioning support