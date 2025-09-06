# Arquitetura - YouMeet

## Visão Geral

O YouMeet segue os princípios da **Arquitetura Hexagonal** (Ports and Adapters), garantindo separação clara de responsabilidades e alta testabilidade.

## Estrutura do Projeto

```
youmeet/
├── cmd/api/                       # Ponto de entrada da aplicação
│   └── main.go
├── internal/
│   ├── core/                      # Núcleo da aplicação (hexagonal)
│   │   ├── domain/               # Entidades e regras de negócio
│   │   │   ├── appointment/      # Domínio de Agendamento
│   │   │   ├── auth/            # Domínio de Autenticação
│   │   │   ├── service/         # Domínio de Serviço
│   │   │   └── user/            # Domínio de Usuário
│   │   └── services/            # Lógica de negócio (casos de uso)
│   │       ├── auth_service.go
│   │       └── booking_service.go
│   ├── adapters/                # Adaptadores externos
│   │   ├── handlers/            # Adaptadores de entrada (HTTP)
│   │   │   ├── auth_handler/
│   │   │   └── appointment_handler/
│   │   └── repositories/        # Adaptadores de saída (dados)
│   │       ├── db_interface.go  # Interface DBClient
│   │       ├── user_repository.go
│   │       ├── appointment_repository.go
│   │       └── service_repository.go
│   └── infra/                   # Implementações de infraestrutura
│       └── database/
│           ├── factory.go       # Factory para escolha do banco
│           ├── postgres_client.go
│           └── sqlite_client.go
└── documentation/               # Documentação do projeto
```

## Princípios da Arquitetura Hexagonal

### 1. **Core (Núcleo)**
- **Domain**: Entidades puras sem dependências externas
- **Services**: Lógica de negócio e casos de uso
- **Não conhece** detalhes de infraestrutura (banco, web, etc.)

### 2. **Adapters (Adaptadores)**
- **Handlers**: Adaptadores de entrada (HTTP → Core)
- **Repositories**: Adaptadores de saída (Core → Database)
- **Interfaces**: Contratos definidos pelos adaptadores

### 3. **Infrastructure (Infraestrutura)**
- **Database Clients**: Implementações concretas dos bancos
- **Factory**: Escolha da implementação baseada em configuração

## Fluxo de Dependências

```
Infrastructure → Adapters → Core
```

- **Core** não depende de nada
- **Adapters** dependem do Core (interfaces)
- **Infrastructure** depende dos Adapters (implementa interfaces)

## Inversão de Dependências

### Interface DBClient
```go
// Definida em: internal/adapters/repositories/db_interface.go
type DBClient interface {
    Create(value interface{}) error
    First(dest interface{}, conds ...interface{}) error
    Find(dest interface{}, conds ...interface{}) error
    Where(query interface{}, args ...interface{}) DBClient
    Delete(value interface{}, conds ...interface{}) error
    AutoMigrate(dst ...interface{}) error
}
```

### Implementações
- **PostgresClient**: Implementa DBClient usando GORM + PostgreSQL
- **SQLiteClient**: Implementa DBClient usando GORM + SQLite

## Benefícios da Arquitetura

### 1. **Testabilidade**
- Core isolado e facilmente testável
- Mocks simples das interfaces
- Testes unitários independentes de infraestrutura

### 2. **Flexibilidade**
- Troca de banco de dados sem alterar lógica de negócio
- Múltiplos adaptadores de entrada (HTTP, CLI, gRPC)
- Fácil adição de novos recursos

### 3. **Manutenibilidade**
- Separação clara de responsabilidades
- Baixo acoplamento entre camadas
- Código limpo e organizad

### 4. **Escalabilidade**
- Arquitetura preparada para crescimento
- Fácil adição de novos domínios
- Suporte a múltiplas implementações

## Padrões Utilizados

### 1. **Repository Pattern**
- Abstração da camada de dados
- Interface definida pelos casos de uso
- Implementações específicas por tecnologia

### 2. **Factory Pattern**
- Criação de objetos baseada em configuração
- Escolha de implementação em tempo de execução
- Facilita testes e configurações diferentes

### 3. **Dependency Injection**
- Injeção de dependências via construtores
- Inversão de controle
- Facilita testes e mocks

## Configuração de Ambiente

### PostgreSQL (Produção)
```bash
export DB_TYPE=postgres
export DATABASE_URL="host=localhost user=youmeet password=youmeet dbname=youmeet port=5432 sslmode=disable"
```

### SQLite (Desenvolvimento/Testes)
```bash
export DB_TYPE=sqlite
export DB_PATH=youmeet.db
```

## Migração de Dados

O sistema utiliza **Auto-Migration** do GORM:
- Tabelas criadas automaticamente a partir das entidades
- Relacionamentos definidos via tags GORM
- Suporte a múltiplos bancos de dados