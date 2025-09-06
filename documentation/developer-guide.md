# Guia do Desenvolvedor - YouMeet

## Configuração do Ambiente de Desenvolvimento

### Pré-requisitos

- **Go 1.21+**
- **Git**
- **PostgreSQL** (opcional, para produção)
- **SQLite** (incluído, para desenvolvimento)

### Instalação

1. **Clone o repositório:**
```bash
git clone https://github.com/Phyper14/youmeet.git
cd youmeet
```

2. **Instale as dependências:**
```bash
go mod download
```

3. **Configure as variáveis de ambiente:**
```bash
# Para desenvolvimento com SQLite
export DB_TYPE=sqlite
export DB_PATH=youmeet.db

# Para produção com PostgreSQL
export DB_TYPE=postgres
export DATABASE_URL="host=localhost user=youmeet password=youmeet dbname=youmeet port=5432 sslmode=disable"
```

4. **Execute a aplicação:**
```bash
go run cmd/api/main.go
```

## Estrutura de Desenvolvimento

### Adicionando um Novo Domínio

1. **Crie a estrutura de domínio:**
```bash
mkdir -p internal/core/domain/novo_dominio
```

2. **Defina as entidades:**
```go
// internal/core/domain/novo_dominio/entidade.go
package novo_dominio

import (
    "time"
    "github.com/google/uuid"
)

type MinhaEntidade struct {
    ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
    Nome      string    `json:"nome" gorm:"not null"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
```

3. **Defina as interfaces (portas):**
```go
// internal/core/domain/novo_dominio/repository.go
package novo_dominio

import (
    "context"
    "github.com/google/uuid"
)

type Repository interface {
    Create(ctx context.Context, entity *MinhaEntidade) error
    GetByID(ctx context.Context, id uuid.UUID) (*MinhaEntidade, error)
}
```

### Implementando um Repository

1. **Crie o repository:**
```go
// internal/adapters/repositories/minha_entidade_repository.go
package repositories

import (
    "context"
    "github.com/google/uuid"
    "youmeet/internal/core/domain/novo_dominio"
)

type MinhaEntidadeRepository struct {
    db DBClient
}

func NewMinhaEntidadeRepository(db DBClient) *MinhaEntidadeRepository {
    return &MinhaEntidadeRepository{db: db}
}

func (r *MinhaEntidadeRepository) Create(ctx context.Context, entity *novo_dominio.MinhaEntidade) error {
    return r.db.Create(entity)
}

func (r *MinhaEntidadeRepository) GetByID(ctx context.Context, id uuid.UUID) (*novo_dominio.MinhaEntidade, error) {
    var entity novo_dominio.MinhaEntidade
    err := r.db.First(&entity, "id = ?", id)
    return &entity, err
}
```

### Criando um Service

1. **Implemente a lógica de negócio:**
```go
// internal/core/services/minha_entidade_service.go
package services

import (
    "context"
    "github.com/google/uuid"
    "youmeet/internal/core/domain/novo_dominio"
)

type MinhaEntidadeService struct {
    repo novo_dominio.Repository
}

func NewMinhaEntidadeService(repo novo_dominio.Repository) *MinhaEntidadeService {
    return &MinhaEntidadeService{repo: repo}
}

func (s *MinhaEntidadeService) CriarEntidade(ctx context.Context, nome string) (*novo_dominio.MinhaEntidade, error) {
    entity := &novo_dominio.MinhaEntidade{
        ID:   uuid.New(),
        Nome: nome,
    }
    
    err := s.repo.Create(ctx, entity)
    return entity, err
}
```

### Adicionando um Handler

1. **Crie a estrutura do handler:**
```bash
mkdir -p internal/adapters/handlers/minha_entidade_handler
```

2. **Implemente o handler:**
```go
// internal/adapters/handlers/minha_entidade_handler/handler.go
package minha_entidade_handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "youmeet/internal/core/services"
)

type Handler struct {
    service *services.MinhaEntidadeService
}

func NewHandler(service *services.MinhaEntidadeService) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
    var req CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    entity, err := h.service.CriarEntidade(c.Request.Context(), req.Nome)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"entity": entity})
}
```

3. **Defina os DTOs:**
```go
// internal/adapters/handlers/minha_entidade_handler/dto.go
package minha_entidade_handler

type CreateRequest struct {
    Nome string `json:"nome" binding:"required"`
}
```

## Padrões de Código

### Nomenclatura

- **Packages**: snake_case (ex: `user_repository`)
- **Types**: PascalCase (ex: `UserRepository`)
- **Functions**: PascalCase para públicas, camelCase para privadas
- **Variables**: camelCase
- **Constants**: UPPER_CASE

### Estrutura de Arquivos

```
domain_name/
├── entity.go      # Entidades do domínio
└── repository.go  # Interfaces (portas)

handler_name/
├── handler.go     # Lógica do handler
└── dto.go         # Data Transfer Objects
```

### Tratamento de Erros

```go
// Sempre retorne erros específicos
func (s *Service) DoSomething() error {
    if err := s.repo.Save(); err != nil {
        return fmt.Errorf("failed to save: %w", err)
    }
    return nil
}

// Use contexto para cancelamento
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Entity, error) {
    // implementação
}
```

## Testes

### Testes Unitários

```go
// internal/core/services/auth_service_test.go
package services_test

import (
    "context"
    "testing"
    "youmeet/internal/core/services"
    "youmeet/internal/mocks"
)

func TestAuthService_Register(t *testing.T) {
    mockRepo := &mocks.UserRepository{}
    service := services.NewAuthService(mockRepo, nil, nil)
    
    // Test implementation
}
```

### Executar Testes

```bash
# Todos os testes
go test ./...

# Testes com coverage
go test -cover ./...

# Testes específicos
go test ./internal/core/services/
```

## Comandos Úteis

### Desenvolvimento

```bash
# Executar aplicação
go run cmd/api/main.go

# Build
go build -o bin/youmeet cmd/api/main.go

# Formatar código
go fmt ./...

# Verificar dependências
go mod tidy

# Verificar código
go vet ./...
```

### Database

```bash
# SQLite (desenvolvimento)
export DB_TYPE=sqlite
export DB_PATH=youmeet.db

# PostgreSQL (produção)
export DB_TYPE=postgres
export DATABASE_URL="postgres://user:pass@localhost/dbname?sslmode=disable"
```

## Contribuindo

1. **Fork** o repositório
2. **Crie** uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. **Commit** suas mudanças (`git commit -am 'Add nova feature'`)
4. **Push** para a branch (`git push origin feature/nova-feature`)
5. **Abra** um Pull Request

### Padrões de Commit

```
feat: adiciona nova funcionalidade
fix: corrige bug
docs: atualiza documentação
style: mudanças de formatação
refactor: refatoração de código
test: adiciona ou modifica testes
chore: mudanças de build ou ferramentas
```