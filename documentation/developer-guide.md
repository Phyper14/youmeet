# Developer Guide

## Development Setup

### Prerequisites

- **Go**: Version 1.21 or higher
- **Git**: For version control
- **IDE**: VS Code, GoLand, or any Go-compatible editor

### Environment Setup

1. **Clone the repository:**
```bash
git clone <repository-url>
cd youmeet
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Verify setup:**
```bash
go build ./cmd/main.go
```

## Development Workflow

### Running the Application

**Development mode:**
```bash
go run cmd/main.go
```

**Build and run:**
```bash
go build -o youmeet cmd/main.go
./youmeet
```

### Testing

**Run all tests:**
```bash
go test ./...
```

**Run tests with coverage:**
```bash
go test -cover ./...
```

**Run tests with verbose output:**
```bash
go test -v ./...
```

### Code Quality

**Format code:**
```bash
go fmt ./...
```

**Lint code:**
```bash
go vet ./...
```

**Static analysis (if golangci-lint is installed):**
```bash
golangci-lint run
```

## Project Structure

```
youmeet/
├── cmd/                    # Application entry points
│   └── main.go            # Main application
├── internal/              # Private application code
│   ├── adapters/          # External adapters
│   │   ├── http_handler.go    # HTTP REST API
│   │   └── memory_repository.go # In-memory storage
│   ├── application/       # Application services
│   │   └── appointment_service.go
│   ├── domain/           # Domain entities
│   │   └── entities.go
│   └── ports/            # Interface definitions
│       ├── repositories.go
│       └── services.go
├── documentation/         # Project documentation
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
└── README.md             # Project overview
```

## Coding Standards

### Go Conventions

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Write clear, concise comments for exported functions
- Keep functions small and focused
- Use interfaces for abstraction

### Architecture Guidelines

1. **Domain Layer**:
   - Keep business logic pure
   - No external dependencies
   - Focus on business rules

2. **Application Layer**:
   - Orchestrate domain entities
   - Implement use cases
   - Depend only on ports

3. **Adapters Layer**:
   - Implement port interfaces
   - Handle external system concerns
   - Convert between formats

4. **Ports Layer**:
   - Define clear interfaces
   - Keep interfaces small and focused
   - Use dependency inversion

### Naming Conventions

- **Packages**: lowercase, single word when possible
- **Files**: snake_case for multi-word names
- **Types**: PascalCase
- **Functions**: PascalCase for exported, camelCase for private
- **Variables**: camelCase
- **Constants**: PascalCase or UPPER_CASE

## Adding New Features

### 1. Define Domain Entity (if needed)

Add new entities to `internal/domain/entities.go`:

```go
type NewEntity struct {
    ID   uuid.UUID `json:"id"`
    Name string    `json:"name"`
}
```

### 2. Define Ports

Add interfaces to appropriate files in `internal/ports/`:

```go
type NewEntityRepository interface {
    Create(ctx context.Context, entity *domain.NewEntity) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.NewEntity, error)
}
```

### 3. Implement Application Service

Create service in `internal/application/`:

```go
type NewEntityService struct {
    repo ports.NewEntityRepository
}

func (s *NewEntityService) CreateEntity(ctx context.Context, name string) (*domain.NewEntity, error) {
    // Implementation
}
```

### 4. Implement Adapters

Add repository implementation to `internal/adapters/`:

```go
func (r *MemoryRepository) CreateNewEntity(ctx context.Context, entity *domain.NewEntity) error {
    // Implementation
}
```

Add HTTP handlers to `internal/adapters/http_handler.go`:

```go
func (h *HTTPHandler) CreateNewEntity(c *gin.Context) {
    // Implementation
}
```

### 5. Wire Dependencies

Update `cmd/main.go` to wire new dependencies:

```go
newEntityService := application.NewNewEntityService(repo)
handler := adapters.NewHTTPHandler(appointmentService, newEntityService)
r.POST("/entities", handler.CreateNewEntity)
```

## Testing Guidelines

### Unit Tests

- Test each layer independently
- Use interfaces for mocking dependencies
- Focus on business logic testing
- Aim for high test coverage

**Example test structure:**
```go
func TestAppointmentService_BookAppointment(t *testing.T) {
    // Arrange
    mockRepo := &MockAppointmentRepository{}
    service := NewAppointmentService(mockRepo, nil)
    
    // Act
    result, err := service.BookAppointment(ctx, serviceID, clientID, startTime)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Integration Tests

- Test adapter implementations
- Test HTTP endpoints
- Use test databases or containers

## Debugging

### Logging

Add structured logging:
```go
import "log"

log.Printf("Booking appointment for client %s", clientID)
```

### Common Issues

1. **Import cycles**: Keep dependencies flowing in one direction
2. **Interface violations**: Ensure adapters implement all interface methods
3. **UUID parsing**: Always validate UUID strings before parsing

## Dependencies

### Current Dependencies

- `github.com/google/uuid` - UUID generation and parsing
- `github.com/gin-gonic/gin` - HTTP web framework

### Adding New Dependencies

1. Add to go.mod:
```bash
go get github.com/new/dependency
```

2. Import in code:
```go
import "github.com/new/dependency"
```

3. Update documentation if it affects API or architecture

## Performance Considerations

- Use context for request cancellation
- Implement proper error handling
- Consider adding database indexes for queries
- Monitor memory usage with in-memory storage
- Add request timeouts for external calls

## Security Best Practices

- Validate all input data
- Use HTTPS in production
- Implement proper authentication
- Sanitize error messages
- Use secure headers
- Implement rate limiting

## Deployment

### Building for Production

```bash
CGO_ENABLED=0 GOOS=linux go build -o youmeet cmd/main.go
```

### Docker (Future)

Create Dockerfile for containerization:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o youmeet cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/youmeet .
CMD ["./youmeet"]
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes following coding standards
4. Add tests for new functionality
5. Update documentation
6. Submit a pull request

## Troubleshooting

### Common Build Issues

- **Module not found**: Run `go mod download`
- **Import cycle**: Check package dependencies
- **Missing dependencies**: Run `go mod tidy`

### Runtime Issues

- **Port already in use**: Change port in main.go or kill existing process
- **Invalid UUID**: Ensure proper UUID format in requests
- **Memory leaks**: Monitor goroutines and memory usage