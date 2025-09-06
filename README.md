# YouMeet

Um sistema de gerenciamento de agendamentos construído em Go, seguindo os princípios da arquitetura hexagonal.

## Visão Geral

YouMeet é uma API RESTful que permite o gerenciamento completo de agendamentos para prestadores de serviços. O sistema suporta três tipos de usuários: clientes, empresas e profissionais autônomos.

## Funcionalidades

### Sistema de Autenticação
- ✅ Registro de usuários (cliente, empresa, profissional)
- ✅ Login com autenticação por senha
- ✅ Gestão de perfis específicos por tipo de usuário

### Gestão de Agendamentos
- ✅ Criação de agendamentos
- ✅ Consulta de agendamentos por cliente
- ✅ Consulta de agendamentos por profissional

### Arquitetura
- ✅ Arquitetura hexagonal (ports and adapters)
- ✅ Separação clara de responsabilidades
- ✅ Fácil testabilidade e manutenibilidade

## Início Rápido

### Pré-requisitos

- Go 1.21 ou superior
- Git

### Instalação

1. Clone o repositório:
```bash
git clone <repository-url>
cd youmeet
```

2. Instale as dependências:
```bash
go mod download
```

3. Execute a aplicação:
```bash
go run cmd/main.go
```

O servidor iniciará em `http://localhost:8080`

## Endpoints da API

### Autenticação
- `POST /auth/register` - Registrar novo usuário
- `POST /auth/login` - Fazer login

### Agendamentos
- `POST /appointments` - Criar agendamento
- `GET /appointments/{clientId}` - Buscar agendamentos do cliente

## Documentação

Documentação completa disponível na pasta `documentation/`:

- 📚 [**API Documentation**](documentation/api.md) - Endpoints, requests e responses
- 🏗️ [**Architecture Guide**](documentation/architecture.md) - Arquitetura hexagonal e princípios
- 👨‍💻 [**Developer Guide**](documentation/developer-guide.md) - Como contribuir e desenvolver
- ⚙️ [**Configuration Guide**](documentation/configuration.md) - Configuração de ambiente e banco
- 🚀 [**Deployment Guide**](documentation/deployment.md) - Deploy local, Docker e nuvem
- 🧪 [**Testing Guide**](documentation/testing.md) - Estratégias e exemplos de testes

## Estrutura do Projeto

```
youmeet/
├── cmd/                           # Ponto de entrada da aplicação
│   └── main.go
├── internal/
│   ├── core/                      # Núcleo da aplicação (hexagonal)
│   │   ├── domain/               # Entidades e regras de negócio
│   │   │   ├── entities.go       # Entidades do domínio
│   │   │   └── ports.go          # Interfaces (portas)
│   │   └── services/             # Lógica de negócio (casos de uso)
│   │       ├── auth_service.go   # Serviço de autenticação
│   │       └── appointment_service.go
│   └── adapters/                 # Adaptadores externos
│       ├── handlers/             # Adaptadores de entrada (HTTP)
│       │   ├── auth_handler.go
│       │   └── appointment_handler.go
│       └── repositories/         # Adaptadores de saída (dados)
│           └── memory_repository.go
├── documentation/                 # Documentação do projeto
└── README.md
```

## Tipos de Usuário

### Cliente (`client`)
- Pode agendar serviços
- Visualizar seus agendamentos

### Empresa (`company`)
- Pode gerenciar funcionários
- Criar e gerenciar serviços
- Visualizar agendamentos da empresa

### Profissional (`professional`)
- Pode ser autônomo ou funcionário de empresa
- Gerenciar disponibilidade
- Visualizar seus agendamentos

## Contribuindo

Leia o [Guia do Desenvolvedor](documentation/developer-guide.md) para instruções de desenvolvimento e contribuição.

## Licença

Este projeto está licenciado sob a Licença MIT.