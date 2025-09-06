# Deploy - YouMeet

## Visão Geral

Este guia cobre diferentes estratégias de deploy para a aplicação YouMeet, desde desenvolvimento local até produção em nuvem.

## Deploy Local

### Desenvolvimento
```bash
# Clone o repositório
git clone https://github.com/Phyper14/youmeet.git
cd youmeet

# Configure ambiente
export DB_TYPE=sqlite
export DB_PATH=youmeet.db

# Execute
go run cmd/api/main.go
```

### Build Local
```bash
# Build para o sistema atual
go build -o bin/youmeet cmd/api/main.go

# Execute o binário
./bin/youmeet
```

### Cross-compilation
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/youmeet-linux cmd/api/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/youmeet.exe cmd/api/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/youmeet-macos cmd/api/main.go
```

## Docker

### Dockerfile
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o youmeet cmd/api/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/

COPY --from=builder /app/youmeet .

EXPOSE 8080

CMD ["./youmeet"]
```

### Build e Run
```bash
# Build da imagem
docker build -t youmeet:latest .

# Run com SQLite
docker run -p 8080:8080 \
  -e DB_TYPE=sqlite \
  -e DB_PATH=/data/youmeet.db \
  -v $(pwd)/data:/data \
  youmeet:latest

# Run com PostgreSQL
docker run -p 8080:8080 \
  -e DB_TYPE=postgres \
  -e DATABASE_URL="postgres://user:pass@host:5432/dbname" \
  youmeet:latest
```

## Docker Compose

### Desenvolvimento
```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_TYPE=sqlite
      - DB_PATH=/data/youmeet.db
      - ENV=development
    volumes:
      - ./data:/data
      - .:/app
    command: go run cmd/api/main.go
```

### Produção
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_TYPE=postgres
      - DATABASE_URL=postgres://youmeet:${DB_PASSWORD}@db:5432/youmeet?sslmode=disable
      - ENV=production
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=youmeet
      - POSTGRES_USER=youmeet
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped

volumes:
  postgres_data:
```

### Comandos
```bash
# Desenvolvimento
docker-compose -f docker-compose.dev.yml up

# Produção
docker-compose -f docker-compose.prod.yml up -d
```

## Deploy em Nuvem

### Heroku

#### Preparação
```bash
# Instale Heroku CLI
# https://devcenter.heroku.com/articles/heroku-cli

# Login
heroku login

# Crie a aplicação
heroku create youmeet-app
```

#### Configuração
```bash
# Configure variáveis de ambiente
heroku config:set DB_TYPE=postgres
heroku config:set ENV=production

# Adicione PostgreSQL
heroku addons:create heroku-postgresql:hobby-dev
```

#### Deploy
```bash
# Deploy via Git
git push heroku main

# Ou via Container Registry
heroku container:push web
heroku container:release web
```

### AWS ECS

#### Task Definition
```json
{
  "family": "youmeet",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "name": "youmeet",
      "image": "youmeet:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "DB_TYPE",
          "value": "postgres"
        },
        {
          "name": "DATABASE_URL",
          "value": "postgres://user:pass@rds-endpoint:5432/youmeet"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/youmeet",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

### Google Cloud Run

#### Preparação
```bash
# Configure gcloud
gcloud auth login
gcloud config set project YOUR_PROJECT_ID

# Build e push para Container Registry
docker build -t gcr.io/YOUR_PROJECT_ID/youmeet .
docker push gcr.io/YOUR_PROJECT_ID/youmeet
```

#### Deploy
```bash
# Deploy para Cloud Run
gcloud run deploy youmeet \
  --image gcr.io/YOUR_PROJECT_ID/youmeet \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars DB_TYPE=postgres,DATABASE_URL="postgres://..."
```

### DigitalOcean App Platform

#### app.yaml
```yaml
name: youmeet
services:
- name: api
  source_dir: /
  github:
    repo: Phyper14/youmeet
    branch: main
  run_command: ./youmeet
  environment_slug: go
  instance_count: 1
  instance_size_slug: basic-xxs
  envs:
  - key: DB_TYPE
    value: postgres
  - key: DATABASE_URL
    value: ${db.DATABASE_URL}

databases:
- name: db
  engine: PG
  version: "13"
  size: db-s-dev-database
```

## Configuração de Proxy Reverso

### Nginx
```nginx
# nginx.conf
upstream youmeet {
    server app:8080;
}

server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://youmeet;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# HTTPS (com SSL)
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;

    location / {
        proxy_pass http://youmeet;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Traefik
```yaml
# docker-compose.yml com Traefik
version: '3.8'

services:
  traefik:
    image: traefik:v2.9
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--entrypoints.web.address=:80"
    ports:
      - "80:80"
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

  app:
    build: .
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.youmeet.rule=Host(`yourdomain.com`)"
      - "traefik.http.services.youmeet.loadbalancer.server.port=8080"
```

## Monitoramento e Logs

### Health Check
```go
// Adicione ao main.go
r.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "healthy",
        "timestamp": time.Now(),
    })
})
```

### Logs Estruturados
```bash
# Configure logs em JSON para produção
export LOG_FORMAT=json
export LOG_LEVEL=info
```

### Prometheus Metrics
```go
// Adicione métricas
import "github.com/prometheus/client_golang/prometheus/promhttp"

r.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

## Backup e Recuperação

### PostgreSQL
```bash
# Backup
pg_dump -h localhost -U youmeet youmeet > backup.sql

# Restore
psql -h localhost -U youmeet youmeet < backup.sql
```

### SQLite
```bash
# Backup
cp youmeet.db backup_$(date +%Y%m%d_%H%M%S).db

# Restore
cp backup_20240115_120000.db youmeet.db
```

## Segurança

### Variáveis de Ambiente Sensíveis
```bash
# Use secrets management
export DATABASE_URL=$(cat /run/secrets/database_url)
export JWT_SECRET=$(cat /run/secrets/jwt_secret)
```

### HTTPS
```bash
# Let's Encrypt com Certbot
certbot --nginx -d yourdomain.com
```

### Firewall
```bash
# UFW (Ubuntu)
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```