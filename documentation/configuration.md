# Configuration Guide

## Overview

YouMeet currently uses minimal configuration with sensible defaults. This guide covers current configuration options and future extensibility.

## Current Configuration

### Server Configuration

The application currently uses hardcoded configuration in `cmd/main.go`:

```go
log.Fatal(http.ListenAndServe(":8080", r))
```

**Default Settings:**
- **Port**: 8080
- **Host**: localhost (0.0.0.0 when deployed)
- **Protocol**: HTTP

### Storage Configuration

Currently uses in-memory storage with no persistence:

```go
repo := adapters.NewMemoryRepository()
```

**Characteristics:**
- **Type**: In-memory maps
- **Persistence**: None (data lost on restart)
- **Concurrency**: Not thread-safe (single-threaded access)

## Environment Variables

### Recommended Environment Variables

While not currently implemented, these environment variables should be supported:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |
| `HOST` | localhost | Server host |
| `LOG_LEVEL` | info | Logging level (debug, info, warn, error) |
| `DB_TYPE` | memory | Database type (memory, postgres, mysql) |
| `DB_URL` | - | Database connection string |
| `DB_MAX_CONNECTIONS` | 10 | Maximum database connections |
| `API_TIMEOUT` | 30s | API request timeout |

### Example Environment Setup

**Development (.env file):**
```bash
PORT=8080
HOST=localhost
LOG_LEVEL=debug
DB_TYPE=memory
```

**Production:**
```bash
PORT=80
HOST=0.0.0.0
LOG_LEVEL=info
DB_TYPE=postgres
DB_URL=postgres://user:pass@localhost/youmeet
DB_MAX_CONNECTIONS=25
API_TIMEOUT=10s
```

## Configuration Implementation

### Recommended Configuration Structure

Create `internal/config/config.go`:

```go
package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Logging  LoggingConfig
}

type ServerConfig struct {
    Port    int
    Host    string
    Timeout time.Duration
}

type DatabaseConfig struct {
    Type           string
    URL            string
    MaxConnections int
}

type LoggingConfig struct {
    Level string
}

func Load() *Config {
    return &Config{
        Server: ServerConfig{
            Port:    getEnvInt("PORT", 8080),
            Host:    getEnv("HOST", "localhost"),
            Timeout: getEnvDuration("API_TIMEOUT", 30*time.Second),
        },
        Database: DatabaseConfig{
            Type:           getEnv("DB_TYPE", "memory"),
            URL:            getEnv("DB_URL", ""),
            MaxConnections: getEnvInt("DB_MAX_CONNECTIONS", 10),
        },
        Logging: LoggingConfig{
            Level: getEnv("LOG_LEVEL", "info"),
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}
```

### Updated Main Function

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "youmeet/internal/config"
    // ... other imports
)

func main() {
    cfg := config.Load()
    
    // Initialize components with config
    repo := adapters.NewMemoryRepository()
    appointmentService := application.NewAppointmentService(repo, repo)
    handler := adapters.NewHTTPHandler(appointmentService)

    r := mux.NewRouter()
    r.HandleFunc("/appointments", handler.BookAppointment).Methods("POST")
    r.HandleFunc("/appointments/{clientId}", handler.GetAppointments).Methods("GET")

    addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    log.Printf("Server starting on %s", addr)
    log.Fatal(http.ListenAndServe(addr, r))
}
```

## Database Configuration

### Memory Storage (Current)

**Configuration:**
```go
repo := adapters.NewMemoryRepository()
```

**Characteristics:**
- No external dependencies
- Fast access
- No persistence
- Limited scalability

### PostgreSQL (Future)

**Environment Variables:**
```bash
DB_TYPE=postgres
DB_URL=postgres://username:password@localhost:5432/youmeet?sslmode=disable
DB_MAX_CONNECTIONS=25
```

**Connection Example:**
```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

func NewPostgresRepository(cfg DatabaseConfig) (*PostgresRepository, error) {
    db, err := sql.Open("postgres", cfg.URL)
    if err != nil {
        return nil, err
    }
    
    db.SetMaxOpenConns(cfg.MaxConnections)
    db.SetMaxIdleConns(cfg.MaxConnections / 2)
    
    return &PostgresRepository{db: db}, nil
}
```

### MySQL (Future)

**Environment Variables:**
```bash
DB_TYPE=mysql
DB_URL=username:password@tcp(localhost:3306)/youmeet
DB_MAX_CONNECTIONS=25
```

## Logging Configuration

### Current Logging

Uses standard Go log package:
```go
log.Println("Server starting on :8080")
```

### Structured Logging (Recommended)

**Environment Variables:**
```bash
LOG_LEVEL=info
LOG_FORMAT=json
```

**Implementation with logrus:**
```go
import "github.com/sirupsen/logrus"

func setupLogging(cfg LoggingConfig) {
    level, err := logrus.ParseLevel(cfg.Level)
    if err != nil {
        level = logrus.InfoLevel
    }
    
    logrus.SetLevel(level)
    logrus.SetFormatter(&logrus.JSONFormatter{})
}
```

## Security Configuration

### HTTPS Configuration

**Environment Variables:**
```bash
TLS_CERT_FILE=/path/to/cert.pem
TLS_KEY_FILE=/path/to/key.pem
HTTPS_ENABLED=true
```

**Implementation:**
```go
if cfg.HTTPS.Enabled {
    log.Fatal(http.ListenAndServeTLS(addr, cfg.HTTPS.CertFile, cfg.HTTPS.KeyFile, r))
} else {
    log.Fatal(http.ListenAndServe(addr, r))
}
```

### CORS Configuration

**Environment Variables:**
```bash
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE
CORS_ALLOWED_HEADERS=Content-Type,Authorization
```

## Monitoring Configuration

### Health Check Endpoint

**Configuration:**
```bash
HEALTH_CHECK_ENABLED=true
HEALTH_CHECK_PATH=/health
```

**Implementation:**
```go
if cfg.HealthCheck.Enabled {
    r.HandleFunc(cfg.HealthCheck.Path, handler.HealthCheck).Methods("GET")
}
```

### Metrics Configuration

**Environment Variables:**
```bash
METRICS_ENABLED=true
METRICS_PATH=/metrics
METRICS_PORT=9090
```

## Configuration Validation

### Validation Function

```go
func (c *Config) Validate() error {
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Server.Port)
    }
    
    if c.Database.Type == "postgres" && c.Database.URL == "" {
        return fmt.Errorf("database URL required for postgres")
    }
    
    return nil
}
```

### Usage in Main

```go
func main() {
    cfg := config.Load()
    if err := cfg.Validate(); err != nil {
        log.Fatalf("Configuration error: %v", err)
    }
    
    // Continue with application startup
}
```

## Configuration Files

### YAML Configuration (Alternative)

Create `config.yaml`:
```yaml
server:
  port: 8080
  host: localhost
  timeout: 30s

database:
  type: memory
  url: ""
  max_connections: 10

logging:
  level: info
```

### JSON Configuration (Alternative)

Create `config.json`:
```json
{
  "server": {
    "port": 8080,
    "host": "localhost",
    "timeout": "30s"
  },
  "database": {
    "type": "memory",
    "url": "",
    "max_connections": 10
  },
  "logging": {
    "level": "info"
  }
}
```

## Deployment Configurations

### Docker Environment

```dockerfile
ENV PORT=8080
ENV HOST=0.0.0.0
ENV LOG_LEVEL=info
ENV DB_TYPE=postgres
ENV DB_URL=postgres://user:pass@db:5432/youmeet
```

### Kubernetes ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: youmeet-config
data:
  PORT: "8080"
  HOST: "0.0.0.0"
  LOG_LEVEL: "info"
  DB_TYPE: "postgres"
```

## Best Practices

1. **Use environment variables** for deployment-specific settings
2. **Provide sensible defaults** for all configuration options
3. **Validate configuration** at startup
4. **Document all options** with examples
5. **Use structured configuration** objects
6. **Support multiple formats** (env vars, files)
7. **Keep secrets separate** from regular configuration
8. **Version configuration schema** for compatibility