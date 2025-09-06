# Testes - YouMeet

## Visão Geral

O YouMeet utiliza uma estratégia de testes em múltiplas camadas, aproveitando a arquitetura hexagonal para facilitar a criação de testes isolados e confiáveis.

## Estrutura de Testes

```
youmeet/
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   │   └── *_test.go          # Testes de entidades
│   │   └── services/
│   │       └── *_service_test.go  # Testes de lógica de negócio
│   ├── adapters/
│   │   ├── handlers/
│   │   │   └── *_handler_test.go  # Testes de handlers HTTP
│   │   └── repositories/
│   │       └── *_repository_test.go # Testes de repositórios
│   └── infra/
│       └── database/
│           └── *_client_test.go   # Testes de clientes de banco
├── tests/
│   ├── integration/               # Testes de integração
│   ├── e2e/                      # Testes end-to-end
│   └── mocks/                    # Mocks e stubs
└── testdata/                     # Dados de teste
```

## Tipos de Testes

### 1. Testes Unitários

#### Testando Entidades (Domain)
```go
// internal/core/domain/user/user_test.go
package user_test

import (
    "testing"
    "youmeet/internal/core/domain/user"
    "github.com/google/uuid"
)

func TestUser_IsValid(t *testing.T) {
    tests := []struct {
        name    string
        user    user.User
        wantErr bool
    }{
        {
            name: "valid user",
            user: user.User{
                ID:    uuid.New(),
                Name:  "João Silva",
                Email: "joao@email.com",
                Role:  "client",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            user: user.User{
                ID:    uuid.New(),
                Name:  "João Silva",
                Email: "invalid-email",
                Role:  "client",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

#### Testando Services (Lógica de Negócio)
```go
// internal/core/services/auth_service_test.go
package services_test

import (
    "context"
    "testing"
    "youmeet/internal/core/domain/auth"
    "youmeet/internal/core/services"
    "youmeet/tests/mocks"
)

func TestAuthService_Register(t *testing.T) {
    // Arrange
    mockUserRepo := &mocks.UserRepository{}
    mockCompanyRepo := &mocks.CompanyRepository{}
    mockProfRepo := &mocks.ProfessionalRepository{}
    
    service := services.NewAuthService(mockUserRepo, mockCompanyRepo, mockProfRepo)
    
    req := &auth.RegisterRequest{
        Name:     "João Silva",
        Email:    "joao@email.com",
        Password: "senha123",
        Role:     "client",
    }

    // Act
    user, err := service.Register(context.Background(), req)

    // Assert
    if err != nil {
        t.Errorf("AuthService.Register() error = %v", err)
    }
    if user == nil {
        t.Error("AuthService.Register() returned nil user")
    }
    if user.Email != req.Email {
        t.Errorf("Expected email %s, got %s", req.Email, user.Email)
    }
}
```

### 2. Testes de Integração

#### Testando Repositórios com Banco Real
```go
// internal/adapters/repositories/user_repository_integration_test.go
package repositories_test

import (
    "context"
    "testing"
    "youmeet/internal/adapters/repositories"
    "youmeet/internal/core/domain/user"
    "youmeet/internal/infra/database"
    "github.com/google/uuid"
)

func TestUserRepository_Integration(t *testing.T) {
    // Setup banco de teste
    db, err := database.NewSQLiteClient(":memory:")
    if err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }

    // Auto-migrate
    err = db.AutoMigrate(&user.User{})
    if err != nil {
        t.Fatalf("Failed to migrate: %v", err)
    }

    repo := repositories.NewUserRepository(db)

    t.Run("Create and Get User", func(t *testing.T) {
        // Arrange
        testUser := &user.User{
            ID:    uuid.New(),
            Name:  "Test User",
            Email: "test@email.com",
            Role:  "client",
        }

        // Act - Create
        err := repo.Create(context.Background(), testUser)
        if err != nil {
            t.Fatalf("Failed to create user: %v", err)
        }

        // Act - Get
        retrievedUser, err := repo.GetByID(context.Background(), testUser.ID)
        if err != nil {
            t.Fatalf("Failed to get user: %v", err)
        }

        // Assert
        if retrievedUser.Email != testUser.Email {
            t.Errorf("Expected email %s, got %s", testUser.Email, retrievedUser.Email)
        }
    })
}
```

### 3. Testes de Handlers (HTTP)

```go
// internal/adapters/handlers/auth_handler/auth_handler_test.go
package auth_handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "youmeet/internal/adapters/handlers/auth_handler"
    "youmeet/internal/core/services"
    "youmeet/tests/mocks"
    "github.com/gin-gonic/gin"
)

func TestAuthHandler_Register(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    
    mockUserRepo := &mocks.UserRepository{}
    mockCompanyRepo := &mocks.CompanyRepository{}
    mockProfRepo := &mocks.ProfessionalRepository{}
    
    authService := services.NewAuthService(mockUserRepo, mockCompanyRepo, mockProfRepo)
    handler := auth_handler.NewHandler(authService)

    router := gin.New()
    router.POST("/auth/register", handler.Register)

    t.Run("successful registration", func(t *testing.T) {
        // Arrange
        reqBody := auth_handler.RegisterRequest{
            Name:     "João Silva",
            Email:    "joao@email.com",
            Password: "senha123",
            Role:     "client",
        }
        
        jsonBody, _ := json.Marshal(reqBody)
        req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
        req.Header.Set("Content-Type", "application/json")
        
        w := httptest.NewRecorder()

        // Act
        router.ServeHTTP(w, req)

        // Assert
        if w.Code != http.StatusCreated {
            t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
        }

        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        if err != nil {
            t.Fatalf("Failed to unmarshal response: %v", err)
        }

        if response["message"] != "Usuário criado com sucesso" {
            t.Errorf("Unexpected response message: %v", response["message"])
        }
    })
}
```

### 4. Testes End-to-End

```go
// tests/e2e/auth_flow_test.go
package e2e_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"
    "youmeet/cmd/api"
)

func TestAuthFlow_E2E(t *testing.T) {
    // Setup servidor de teste
    server := api.SetupTestServer()
    defer server.Close()

    t.Run("complete auth flow", func(t *testing.T) {
        // 1. Register
        registerReq := map[string]string{
            "name":     "Test User",
            "email":    "test@email.com",
            "password": "senha123",
            "role":     "client",
        }
        
        registerBody, _ := json.Marshal(registerReq)
        resp, err := http.Post(server.URL+"/auth/register", "application/json", bytes.NewBuffer(registerBody))
        if err != nil {
            t.Fatalf("Register request failed: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusCreated {
            t.Errorf("Expected status 201, got %d", resp.StatusCode)
        }

        // 2. Login
        loginReq := map[string]string{
            "email":    "test@email.com",
            "password": "senha123",
        }
        
        loginBody, _ := json.Marshal(loginReq)
        resp, err = http.Post(server.URL+"/auth/login", "application/json", bytes.NewBuffer(loginBody))
        if err != nil {
            t.Fatalf("Login request failed: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Errorf("Expected status 200, got %d", resp.StatusCode)
        }

        var loginResponse map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&loginResponse)
        
        if loginResponse["token"] == nil {
            t.Error("Expected token in login response")
        }
    })
}
```

## Mocks e Stubs

### Criando Mocks
```go
// tests/mocks/user_repository.go
package mocks

import (
    "context"
    "youmeet/internal/core/domain/user"
    "github.com/google/uuid"
)

type UserRepository struct {
    users map[uuid.UUID]*user.User
}

func NewUserRepository() *UserRepository {
    return &UserRepository{
        users: make(map[uuid.UUID]*user.User),
    }
}

func (m *UserRepository) Create(ctx context.Context, u *user.User) error {
    m.users[u.ID] = u
    return nil
}

func (m *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
    if user, exists := m.users[id]; exists {
        return user, nil
    }
    return nil, errors.New("user not found")
}

func (m *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
    for _, user := range m.users {
        if user.Email == email {
            return user, nil
        }
    }
    return nil, errors.New("user not found")
}
```

## Configuração de Testes

### Makefile
```makefile
# Makefile
.PHONY: test test-unit test-integration test-e2e test-coverage

test: test-unit test-integration

test-unit:
	go test -short ./internal/core/...

test-integration:
	go test -tags=integration ./internal/adapters/...

test-e2e:
	go test -tags=e2e ./tests/e2e/...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-watch:
	find . -name "*.go" | entr -r make test-unit
```

### GitHub Actions
```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: youmeet_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run unit tests
      run: make test-unit
    
    - name: Run integration tests
      run: make test-integration
      env:
        DB_TYPE: postgres
        DATABASE_URL: postgres://postgres:postgres@localhost:5432/youmeet_test?sslmode=disable
    
    - name: Generate coverage report
      run: make test-coverage
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
```

## Comandos de Teste

### Executar Testes
```bash
# Todos os testes
go test ./...

# Apenas testes unitários
go test -short ./...

# Testes com coverage
go test -cover ./...

# Testes específicos
go test ./internal/core/services/

# Testes com verbose
go test -v ./...

# Testes em paralelo
go test -parallel 4 ./...
```

### Benchmarks
```go
func BenchmarkAuthService_Register(b *testing.B) {
    service := setupAuthService()
    req := &auth.RegisterRequest{
        Name:     "Test User",
        Email:    "test@email.com",
        Password: "senha123",
        Role:     "client",
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.Register(context.Background(), req)
    }
}
```

```bash
# Executar benchmarks
go test -bench=. ./...

# Benchmark com memory profiling
go test -bench=. -benchmem ./...
```

## Boas Práticas

### 1. **Arrange, Act, Assert (AAA)**
```go
func TestSomething(t *testing.T) {
    // Arrange - Setup
    input := "test"
    expected := "expected"
    
    // Act - Execute
    result := DoSomething(input)
    
    // Assert - Verify
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

### 2. **Table-Driven Tests**
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@email.com", false},
        {"invalid email", "invalid", true},
        {"empty email", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 3. **Test Helpers**
```go
// tests/helpers/test_helpers.go
package helpers

import (
    "testing"
    "youmeet/internal/core/domain/user"
    "github.com/google/uuid"
)

func CreateTestUser(t *testing.T) *user.User {
    t.Helper()
    return &user.User{
        ID:    uuid.New(),
        Name:  "Test User",
        Email: "test@email.com",
        Role:  "client",
    }
}
```

### 4. **Cleanup**
```go
func TestWithCleanup(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    
    // Cleanup
    t.Cleanup(func() {
        db.Close()
    })
    
    // Test logic
}
```