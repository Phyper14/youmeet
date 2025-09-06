# Deployment Guide

## Overview

This guide covers various deployment options for the YouMeet application, from local development to production environments.

## Prerequisites

- Go 1.21 or higher
- Git
- Docker (for containerized deployments)
- Access to target deployment environment

## Local Deployment

### Development Server

**Quick start:**
```bash
go run cmd/main.go
```

**Build and run:**
```bash
go build -o youmeet cmd/main.go
./youmeet
```

**With custom port:**
```bash
PORT=3000 go run cmd/main.go
```

### Production Build

**Linux/macOS:**
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o youmeet cmd/main.go
```

**Windows:**
```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o youmeet.exe cmd/main.go
```

**Cross-platform builds:**
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o youmeet-linux cmd/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o youmeet-macos cmd/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o youmeet-windows.exe cmd/main.go
```

## Docker Deployment

### Dockerfile

Create `Dockerfile` in project root:
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o youmeet cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/youmeet .

EXPOSE 8080

CMD ["./youmeet"]
```

### Docker Commands

**Build image:**
```bash
docker build -t youmeet:latest .
```

**Run container:**
```bash
docker run -p 8080:8080 youmeet:latest
```

**Run with environment variables:**
```bash
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e LOG_LEVEL=info \
  youmeet:latest
```

### Docker Compose

Create `docker-compose.yml`:
```yaml
version: '3.8'

services:
  youmeet:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - LOG_LEVEL=info
    restart: unless-stopped

  # Future: Add database service
  # postgres:
  #   image: postgres:15
  #   environment:
  #     POSTGRES_DB: youmeet
  #     POSTGRES_USER: youmeet
  #     POSTGRES_PASSWORD: password
  #   volumes:
  #     - postgres_data:/var/lib/postgresql/data
  #   ports:
  #     - "5432:5432"

# volumes:
#   postgres_data:
```

**Run with Docker Compose:**
```bash
docker-compose up -d
```

## Cloud Deployments

### AWS Deployment

#### AWS ECS (Elastic Container Service)

**1. Build and push to ECR:**
```bash
# Create ECR repository
aws ecr create-repository --repository-name youmeet

# Get login token
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

# Tag and push image
docker tag youmeet:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/youmeet:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/youmeet:latest
```

**2. Create ECS task definition:**
```json
{
  "family": "youmeet",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::<account-id>:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "name": "youmeet",
      "image": "<account-id>.dkr.ecr.us-east-1.amazonaws.com/youmeet:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "PORT",
          "value": "8080"
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

#### AWS Lambda (Serverless)

**1. Install AWS Lambda Go runtime:**
```bash
go get github.com/aws/aws-lambda-go/lambda
```

**2. Create Lambda handler:**
```go
package main

import (
    "context"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Initialize your router here
    r := setupRouter()
    adapter := gorillamux.New(r)
    return adapter.ProxyWithContext(ctx, req)
}

func main() {
    lambda.Start(handler)
}
```

### Google Cloud Platform

#### Cloud Run

**1. Build and deploy:**
```bash
gcloud builds submit --tag gcr.io/PROJECT-ID/youmeet
gcloud run deploy --image gcr.io/PROJECT-ID/youmeet --platform managed
```

**2. Using Cloud Build:**

Create `cloudbuild.yaml`:
```yaml
steps:
- name: 'gcr.io/cloud-builders/docker'
  args: ['build', '-t', 'gcr.io/$PROJECT_ID/youmeet', '.']
- name: 'gcr.io/cloud-builders/docker'
  args: ['push', 'gcr.io/$PROJECT_ID/youmeet']
- name: 'gcr.io/cloud-builders/gcloud'
  args: ['run', 'deploy', 'youmeet', '--image', 'gcr.io/$PROJECT_ID/youmeet', '--region', 'us-central1', '--platform', 'managed']
```

### Microsoft Azure

#### Azure Container Instances

```bash
az container create \
  --resource-group myResourceGroup \
  --name youmeet \
  --image youmeet:latest \
  --dns-name-label youmeet-app \
  --ports 8080
```

#### Azure App Service

```bash
az webapp create \
  --resource-group myResourceGroup \
  --plan myAppServicePlan \
  --name youmeet-app \
  --deployment-container-image-name youmeet:latest
```

## Kubernetes Deployment

### Deployment Manifest

Create `k8s/deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: youmeet
  labels:
    app: youmeet
spec:
  replicas: 3
  selector:
    matchLabels:
      app: youmeet
  template:
    metadata:
      labels:
        app: youmeet
    spec:
      containers:
      - name: youmeet
        image: youmeet:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Service Manifest

Create `k8s/service.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: youmeet-service
spec:
  selector:
    app: youmeet
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

### ConfigMap

Create `k8s/configmap.yaml`:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: youmeet-config
data:
  PORT: "8080"
  LOG_LEVEL: "info"
  DB_TYPE: "memory"
```

### Deploy to Kubernetes

```bash
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

## Monitoring and Logging

### Health Check Endpoint

Add to your application:
```go
func (h *HTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
    })
}
```

### Prometheus Metrics

**1. Add Prometheus dependency:**
```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

**2. Add metrics endpoint:**
```go
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    // ... existing code
    
    r.Handle("/metrics", promhttp.Handler())
    
    // ... rest of main function
}
```

### Structured Logging

**Using logrus:**
```bash
go get github.com/sirupsen/logrus
```

```go
import "github.com/sirupsen/logrus"

func init() {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.SetLevel(logrus.InfoLevel)
}
```

## Environment-Specific Configurations

### Development

```bash
export PORT=8080
export LOG_LEVEL=debug
export DB_TYPE=memory
```

### Staging

```bash
export PORT=8080
export LOG_LEVEL=info
export DB_TYPE=postgres
export DB_URL=postgres://user:pass@staging-db:5432/youmeet
```

### Production

```bash
export PORT=80
export LOG_LEVEL=warn
export DB_TYPE=postgres
export DB_URL=postgres://user:pass@prod-db:5432/youmeet
export DB_MAX_CONNECTIONS=25
```

## Security Considerations

### HTTPS/TLS

**Using Let's Encrypt with reverse proxy:**
```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;
    
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Firewall Rules

**Allow only necessary ports:**
```bash
# Allow HTTP/HTTPS
sudo ufw allow 80
sudo ufw allow 443

# Allow SSH (if needed)
sudo ufw allow 22

# Block application port from external access
sudo ufw deny 8080
```

## Backup and Recovery

### Database Backup (Future)

```bash
# PostgreSQL backup
pg_dump -h localhost -U youmeet youmeet > backup.sql

# Restore
psql -h localhost -U youmeet youmeet < backup.sql
```

### Application State

Since the current implementation uses in-memory storage, consider:
- Implementing persistent storage
- Regular data exports
- State replication across instances

## Troubleshooting

### Common Issues

**Port already in use:**
```bash
# Find process using port
lsof -i :8080
# Kill process
kill -9 <PID>
```

**Memory issues:**
```bash
# Check memory usage
docker stats
# Increase container memory limits
```

**Container won't start:**
```bash
# Check logs
docker logs <container-id>
# Check container status
docker ps -a
```

### Debugging

**Enable debug logging:**
```bash
LOG_LEVEL=debug ./youmeet
```

**Check application health:**
```bash
curl http://localhost:8080/health
```

**Monitor resource usage:**
```bash
# CPU and memory
top
# Network connections
netstat -tulpn
```