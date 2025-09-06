# Configuração - YouMeet

## Variáveis de Ambiente

### Configuração de Banco de Dados

#### PostgreSQL (Produção)
```bash
# Tipo do banco de dados
DB_TYPE=postgres

# String de conexão completa
DATABASE_URL="host=localhost user=youmeet password=youmeet dbname=youmeet port=5432 sslmode=disable"

# Ou configuração individual
DB_HOST=localhost
DB_PORT=5432
DB_USER=youmeet
DB_PASSWORD=youmeet
DB_NAME=youmeet
DB_SSLMODE=disable
```

#### SQLite (Desenvolvimento/Testes)
```bash
# Tipo do banco de dados
DB_TYPE=sqlite

# Caminho do arquivo de banco
DB_PATH=youmeet.db

# Para testes em memória
DB_PATH=:memory:
```

### Configuração do Servidor
```bash
# Porta do servidor (padrão: 8080)
PORT=8080

# Ambiente de execução
ENV=development  # development, staging, production

# Nível de log
LOG_LEVEL=info   # debug, info, warn, error
```

## Arquivos de Configuração

### .env (Desenvolvimento)
```bash
# Database
DB_TYPE=sqlite
DB_PATH=youmeet.db

# Server
PORT=8080
ENV=development
LOG_LEVEL=debug
```

### .env.production (Produção)
```bash
# Database
DB_TYPE=postgres
DATABASE_URL=postgres://user:password@localhost:5432/youmeet?sslmode=disable

# Server
PORT=8080
ENV=production
LOG_LEVEL=info
```

## Configuração por Ambiente

### Desenvolvimento Local

1. **Copie o arquivo de exemplo:**
```bash
cp .env.example .env
```

2. **Configure para SQLite:**
```bash
DB_TYPE=sqlite
DB_PATH=youmeet.db
PORT=8080
```

3. **Execute a aplicação:**
```bash
go run cmd/api/main.go
```

### Testes

1. **Configure variáveis para teste:**
```bash
export DB_TYPE=sqlite
export DB_PATH=:memory:
export ENV=test
```

2. **Execute os testes:**
```bash
go test ./...
```

### Produção

#### Docker
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o youmeet cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/youmeet .
CMD ["./youmeet"]
```

#### Docker Compose
```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_TYPE=postgres
      - DATABASE_URL=postgres://youmeet:youmeet@db:5432/youmeet?sslmode=disable
    depends_on:
      - db

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=youmeet
      - POSTGRES_USER=youmeet
      - POSTGRES_PASSWORD=youmeet
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## Configuração de Banco de Dados

### PostgreSQL

#### Instalação Local
```bash
# Ubuntu/Debian
sudo apt-get install postgresql postgresql-contrib

# macOS
brew install postgresql

# Windows
# Baixe do site oficial: https://www.postgresql.org/download/windows/
```

#### Configuração
```sql
-- Conecte como superuser
sudo -u postgres psql

-- Crie o banco e usuário
CREATE DATABASE youmeet;
CREATE USER youmeet WITH PASSWORD 'youmeet';
GRANT ALL PRIVILEGES ON DATABASE youmeet TO youmeet;
```

#### String de Conexão
```bash
DATABASE_URL="postgres://youmeet:youmeet@localhost:5432/youmeet?sslmode=disable"
```

### SQLite

#### Vantagens para Desenvolvimento
- **Sem instalação**: Incluído no Go
- **Arquivo único**: Fácil backup e reset
- **Performance**: Rápido para desenvolvimento
- **Portabilidade**: Funciona em qualquer OS

#### Configuração
```bash
# Arquivo local
DB_PATH=youmeet.db

# Em memória (testes)
DB_PATH=:memory:
```

## Configuração de Logs

### Níveis de Log
- **debug**: Informações detalhadas para debugging
- **info**: Informações gerais sobre operações
- **warn**: Avisos sobre situações não ideais
- **error**: Erros que precisam de atenção

### Configuração
```bash
# Desenvolvimento
LOG_LEVEL=debug

# Produção
LOG_LEVEL=info
```

## Configuração de CORS

Para desenvolvimento frontend:

```go
// cmd/api/main.go
import "github.com/gin-contrib/cors"

func main() {
    r := gin.Default()
    
    // CORS middleware
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
    
    // ... resto da configuração
}
```

## Configuração de Segurança

### JWT (Futuro)
```bash
# Chave secreta para JWT
JWT_SECRET=sua-chave-secreta-muito-segura

# Tempo de expiração
JWT_EXPIRATION=24h
```

### Rate Limiting
```bash
# Requests por minuto por IP
RATE_LIMIT=100

# Burst size
RATE_BURST=10
```

## Monitoramento

### Health Check
```bash
# Endpoint de saúde
curl http://localhost:8080/health
```

### Métricas
```bash
# Endpoint de métricas (Prometheus)
curl http://localhost:8080/metrics
```

## Troubleshooting

### Problemas Comuns

#### Erro de Conexão com Banco
```bash
# Verifique se o PostgreSQL está rodando
sudo systemctl status postgresql

# Teste a conexão
psql -h localhost -U youmeet -d youmeet
```

#### Porta em Uso
```bash
# Encontre o processo usando a porta
lsof -i :8080

# Mate o processo
kill -9 <PID>
```

#### Permissões de Arquivo (SQLite)
```bash
# Verifique permissões
ls -la youmeet.db

# Corrija permissões
chmod 664 youmeet.db
```

### Logs de Debug

```bash
# Execute com logs detalhados
LOG_LEVEL=debug go run cmd/api/main.go

# Ou com variável de ambiente
export LOG_LEVEL=debug
go run cmd/api/main.go
```