# Testing Guide

## Overview

This guide covers testing strategies, implementation patterns, and best practices for the YouMeet application.

## Testing Strategy

### Testing Pyramid

1. **Unit Tests** (70%) - Test individual components in isolation
2. **Integration Tests** (20%) - Test component interactions
3. **End-to-End Tests** (10%) - Test complete user workflows

### Test Categories

- **Domain Tests** - Business logic validation
- **Application Tests** - Use case testing with mocked dependencies
- **Adapter Tests** - External system integration testing
- **API Tests** - HTTP endpoint testing

## Unit Testing

### Domain Layer Testing

**Test file: `internal/domain/entities_test.go`**
```go
package domain

import (
    "testing"
    "time"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)

func TestAppointment_Creation(t *testing.T) {
    // Arrange
    id := uuid.New()
    serviceID := uuid.New()
    clientID := uuid.New()
    startTime := time.Now()

    // Act
    appointment := &Appointment{
        ID:        id,
        ServiceID: serviceID,
        ClientID:  clientID,
        StartTime: startTime,
        Status:    "scheduled",
    }

    // Assert
    assert.Equal(t, id, appointment.ID)
    assert.Equal(t, serviceID, appointment.ServiceID)
    assert.Equal(t, clientID, appointment.ClientID)
    assert.Equal(t, startTime, appointment.StartTime)
    assert.Equal(t, "scheduled", appointment.Status)
}

func TestService_Validation(t *testing.T) {
    tests := []struct {
        name     string
        service  Service
        expected bool
    }{
        {
            name: "valid service",
            service: Service{
                ID:       uuid.New(),
                Name:     "Consultation",
                Duration: 30 * time.Minute,
                Price:    100.0,
            },
            expected: true,
        },
        {
            name: "invalid price",
            service: Service{
                ID:       uuid.New(),
                Name:     "Consultation",
                Duration: 30 * time.Minute,
                Price:    -10.0,
            },
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            valid := tt.service.Price > 0 && tt.service.Name != ""
            assert.Equal(t, tt.expected, valid)
        })
    }
}
```

### Application Layer Testing

**Test file: `internal/application/appointment_service_test.go`**
```go
package application

import (
    "context"
    "testing"
    "time"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "youmeet/internal/domain"
)

// Mock repository
type MockAppointmentRepository struct {
    mock.Mock
}

func (m *MockAppointmentRepository) Create(ctx context.Context, appointment *domain.Appointment) error {
    args := m.Called(ctx, appointment)
    return args.Error(0)
}

func (m *MockAppointmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) List(ctx context.Context, clientID uuid.UUID) ([]*domain.Appointment, error) {
    args := m.Called(ctx, clientID)
    return args.Get(0).([]*domain.Appointment), args.Error(1)
}

type MockServiceRepository struct {
    mock.Mock
}

func (m *MockServiceRepository) Create(ctx context.Context, service *domain.Service) error {
    args := m.Called(ctx, service)
    return args.Error(0)
}

func (m *MockServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*domain.Service), args.Error(1)
}

func (m *MockServiceRepository) List(ctx context.Context) ([]*domain.Service, error) {
    args := m.Called(ctx)
    return args.Get(0).([]*domain.Service), args.Error(1)
}

func TestAppointmentService_BookAppointment(t *testing.T) {
    // Arrange
    mockAppointmentRepo := new(MockAppointmentRepository)
    mockServiceRepo := new(MockServiceRepository)
    service := NewAppointmentService(mockAppointmentRepo, mockServiceRepo)
    
    ctx := context.Background()
    serviceID := uuid.New()
    clientID := uuid.New()
    startTimeStr := "2024-01-15T10:00:00Z"
    
    mockAppointmentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Appointment")).Return(nil)

    // Act
    appointment, err := service.BookAppointment(ctx, serviceID, clientID, startTimeStr)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, appointment)
    assert.Equal(t, serviceID, appointment.ServiceID)
    assert.Equal(t, clientID, appointment.ClientID)
    assert.Equal(t, "scheduled", appointment.Status)
    mockAppointmentRepo.AssertExpectations(t)
}

func TestAppointmentService_BookAppointment_InvalidTime(t *testing.T) {
    // Arrange
    mockAppointmentRepo := new(MockAppointmentRepository)
    mockServiceRepo := new(MockServiceRepository)
    service := NewAppointmentService(mockAppointmentRepo, mockServiceRepo)
    
    ctx := context.Background()
    serviceID := uuid.New()
    clientID := uuid.New()
    invalidTimeStr := "invalid-time-format"

    // Act
    appointment, err := service.BookAppointment(ctx, serviceID, clientID, invalidTimeStr)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, appointment)
}

func TestAppointmentService_GetAppointments(t *testing.T) {
    // Arrange
    mockAppointmentRepo := new(MockAppointmentRepository)
    mockServiceRepo := new(MockServiceRepository)
    service := NewAppointmentService(mockAppointmentRepo, mockServiceRepo)
    
    ctx := context.Background()
    clientID := uuid.New()
    
    expectedAppointments := []*domain.Appointment{
        {
            ID:        uuid.New(),
            ServiceID: uuid.New(),
            ClientID:  clientID,
            StartTime: time.Now(),
            Status:    "scheduled",
        },
    }
    
    mockAppointmentRepo.On("List", ctx, clientID).Return(expectedAppointments, nil)

    // Act
    appointments, err := service.GetAppointments(ctx, clientID)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedAppointments, appointments)
    mockAppointmentRepo.AssertExpectations(t)
}
```

### Adapter Layer Testing

**Test file: `internal/adapters/memory_repository_test.go`**
```go
package adapters

import (
    "context"
    "testing"
    "time"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "youmeet/internal/domain"
)

func TestMemoryRepository_CreateAndGetAppointment(t *testing.T) {
    // Arrange
    repo := NewMemoryRepository()
    ctx := context.Background()
    
    appointment := &domain.Appointment{
        ID:        uuid.New(),
        ServiceID: uuid.New(),
        ClientID:  uuid.New(),
        StartTime: time.Now(),
        Status:    "scheduled",
    }

    // Act - Create
    err := repo.Create(ctx, appointment)
    assert.NoError(t, err)

    // Act - Get
    retrieved, err := repo.GetByID(ctx, appointment.ID)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, appointment, retrieved)
}

func TestMemoryRepository_ListAppointmentsByClient(t *testing.T) {
    // Arrange
    repo := NewMemoryRepository()
    ctx := context.Background()
    clientID := uuid.New()
    
    appointment1 := &domain.Appointment{
        ID:        uuid.New(),
        ServiceID: uuid.New(),
        ClientID:  clientID,
        StartTime: time.Now(),
        Status:    "scheduled",
    }
    
    appointment2 := &domain.Appointment{
        ID:        uuid.New(),
        ServiceID: uuid.New(),
        ClientID:  clientID,
        StartTime: time.Now().Add(time.Hour),
        Status:    "scheduled",
    }
    
    // Different client
    appointment3 := &domain.Appointment{
        ID:        uuid.New(),
        ServiceID: uuid.New(),
        ClientID:  uuid.New(),
        StartTime: time.Now().Add(2 * time.Hour),
        Status:    "scheduled",
    }

    // Act
    repo.Create(ctx, appointment1)
    repo.Create(ctx, appointment2)
    repo.Create(ctx, appointment3)
    
    appointments, err := repo.List(ctx, clientID)

    // Assert
    assert.NoError(t, err)
    assert.Len(t, appointments, 2)
    assert.Contains(t, appointments, appointment1)
    assert.Contains(t, appointments, appointment2)
}

func TestMemoryRepository_GetNonExistentAppointment(t *testing.T) {
    // Arrange
    repo := NewMemoryRepository()
    ctx := context.Background()
    nonExistentID := uuid.New()

    // Act
    appointment, err := repo.GetByID(ctx, nonExistentID)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, appointment)
    assert.Contains(t, err.Error(), "not found")
}
```

## Integration Testing

### HTTP Handler Testing

**Test file: `internal/adapters/http_handler_test.go`**
```go
package adapters

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gorilla/mux"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "youmeet/internal/domain"
)

type MockAppointmentService struct {
    mock.Mock
}

func (m *MockAppointmentService) BookAppointment(ctx context.Context, serviceID, clientID uuid.UUID, startTime string) (*domain.Appointment, error) {
    args := m.Called(ctx, serviceID, clientID, startTime)
    return args.Get(0).(*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentService) GetAppointments(ctx context.Context, clientID uuid.UUID) ([]*domain.Appointment, error) {
    args := m.Called(ctx, clientID)
    return args.Get(0).([]*domain.Appointment), args.Error(1)
}

func TestHTTPHandler_BookAppointment(t *testing.T) {
    // Arrange
    mockService := new(MockAppointmentService)
    handler := NewHTTPHandler(mockService)
    
    serviceID := uuid.New()
    clientID := uuid.New()
    startTime := "2024-01-15T10:00:00Z"
    
    expectedAppointment := &domain.Appointment{
        ID:        uuid.New(),
        ServiceID: serviceID,
        ClientID:  clientID,
        Status:    "scheduled",
    }
    
    mockService.On("BookAppointment", mock.Anything, serviceID, clientID, startTime).Return(expectedAppointment, nil)
    
    requestBody := BookAppointmentRequest{
        ServiceID: serviceID,
        ClientID:  clientID,
        StartTime: startTime,
    }
    
    jsonBody, _ := json.Marshal(requestBody)
    req := httptest.NewRequest("POST", "/appointments", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    
    rr := httptest.NewRecorder()

    // Act
    handler.BookAppointment(rr, req)

    // Assert
    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
    
    var response domain.Appointment
    err := json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, expectedAppointment.ID, response.ID)
    
    mockService.AssertExpectations(t)
}

func TestHTTPHandler_BookAppointment_InvalidJSON(t *testing.T) {
    // Arrange
    mockService := new(MockAppointmentService)
    handler := NewHTTPHandler(mockService)
    
    req := httptest.NewRequest("POST", "/appointments", bytes.NewBufferString("invalid json"))
    req.Header.Set("Content-Type", "application/json")
    
    rr := httptest.NewRecorder()

    // Act
    handler.BookAppointment(rr, req)

    // Assert
    assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHTTPHandler_GetAppointments(t *testing.T) {
    // Arrange
    mockService := new(MockAppointmentService)
    handler := NewHTTPHandler(mockService)
    
    clientID := uuid.New()
    expectedAppointments := []*domain.Appointment{
        {
            ID:        uuid.New(),
            ServiceID: uuid.New(),
            ClientID:  clientID,
            Status:    "scheduled",
        },
    }
    
    mockService.On("GetAppointments", mock.Anything, clientID).Return(expectedAppointments, nil)
    
    req := httptest.NewRequest("GET", "/appointments/"+clientID.String(), nil)
    req = mux.SetURLVars(req, map[string]string{"clientId": clientID.String()})
    
    rr := httptest.NewRecorder()

    // Act
    handler.GetAppointments(rr, req)

    // Assert
    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
    
    var response []*domain.Appointment
    err := json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Len(t, response, 1)
    assert.Equal(t, expectedAppointments[0].ID, response[0].ID)
    
    mockService.AssertExpectations(t)
}

func TestHTTPHandler_GetAppointments_InvalidClientID(t *testing.T) {
    // Arrange
    mockService := new(MockAppointmentService)
    handler := NewHTTPHandler(mockService)
    
    req := httptest.NewRequest("GET", "/appointments/invalid-uuid", nil)
    req = mux.SetURLVars(req, map[string]string{"clientId": "invalid-uuid"})
    
    rr := httptest.NewRecorder()

    // Act
    handler.GetAppointments(rr, req)

    // Assert
    assert.Equal(t, http.StatusBadRequest, rr.Code)
}
```

## End-to-End Testing

### API Integration Tests

**Test file: `test/integration/api_test.go`**
```go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gorilla/mux"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "youmeet/internal/adapters"
    "youmeet/internal/application"
)

func setupTestServer() *httptest.Server {
    repo := adapters.NewMemoryRepository()
    appointmentService := application.NewAppointmentService(repo, repo)
    handler := adapters.NewHTTPHandler(appointmentService)

    r := mux.NewRouter()
    r.HandleFunc("/appointments", handler.BookAppointment).Methods("POST")
    r.HandleFunc("/appointments/{clientId}", handler.GetAppointments).Methods("GET")

    return httptest.NewServer(r)
}

func TestAPI_BookAndRetrieveAppointment(t *testing.T) {
    // Arrange
    server := setupTestServer()
    defer server.Close()

    serviceID := uuid.New()
    clientID := uuid.New()
    startTime := "2024-01-15T10:00:00Z"

    bookingRequest := map[string]interface{}{
        "service_id": serviceID,
        "client_id":  clientID,
        "start_time": startTime,
    }

    jsonBody, _ := json.Marshal(bookingRequest)

    // Act 1: Book appointment
    resp, err := http.Post(server.URL+"/appointments", "application/json", bytes.NewBuffer(jsonBody))
    assert.NoError(t, err)
    defer resp.Body.Close()

    // Assert 1: Booking successful
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var bookedAppointment map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&bookedAppointment)
    assert.NoError(t, err)
    assert.Equal(t, serviceID.String(), bookedAppointment["service_id"])
    assert.Equal(t, clientID.String(), bookedAppointment["client_id"])

    // Act 2: Retrieve appointments
    resp2, err := http.Get(server.URL + "/appointments/" + clientID.String())
    assert.NoError(t, err)
    defer resp2.Body.Close()

    // Assert 2: Retrieval successful
    assert.Equal(t, http.StatusOK, resp2.StatusCode)

    var appointments []map[string]interface{}
    err = json.NewDecoder(resp2.Body).Decode(&appointments)
    assert.NoError(t, err)
    assert.Len(t, appointments, 1)
    assert.Equal(t, serviceID.String(), appointments[0]["service_id"])
}

func TestAPI_MultipleAppointmentsForClient(t *testing.T) {
    // Arrange
    server := setupTestServer()
    defer server.Close()

    clientID := uuid.New()
    
    appointments := []map[string]interface{}{
        {
            "service_id": uuid.New(),
            "client_id":  clientID,
            "start_time": "2024-01-15T10:00:00Z",
        },
        {
            "service_id": uuid.New(),
            "client_id":  clientID,
            "start_time": "2024-01-15T11:00:00Z",
        },
    }

    // Act: Book multiple appointments
    for _, appointment := range appointments {
        jsonBody, _ := json.Marshal(appointment)
        resp, err := http.Post(server.URL+"/appointments", "application/json", bytes.NewBuffer(jsonBody))
        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)
        resp.Body.Close()
    }

    // Act: Retrieve all appointments
    resp, err := http.Get(server.URL + "/appointments/" + clientID.String())
    assert.NoError(t, err)
    defer resp.Body.Close()

    // Assert
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var retrievedAppointments []map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&retrievedAppointments)
    assert.NoError(t, err)
    assert.Len(t, retrievedAppointments, 2)
}
```

## Test Configuration

### Test Dependencies

Add to `go.mod`:
```go
require (
    github.com/stretchr/testify v1.8.4
)
```

Install test dependencies:
```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/stretchr/testify/suite
```

### Test Makefile

Create `Makefile`:
```makefile
.PHONY: test test-unit test-integration test-coverage test-race

# Run all tests
test:
	go test ./...

# Run unit tests only
test-unit:
	go test -short ./...

# Run integration tests only
test-integration:
	go test -run Integration ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
test-race:
	go test -race ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Clean test cache
test-clean:
	go clean -testcache
```

## Test Execution

### Running Tests

**All tests:**
```bash
go test ./...
```

**Specific package:**
```bash
go test ./internal/application
```

**With coverage:**
```bash
go test -cover ./...
```

**With race detection:**
```bash
go test -race ./...
```

**Verbose output:**
```bash
go test -v ./...
```

**Short tests only (skip integration):**
```bash
go test -short ./...
```

### Test Coverage

**Generate coverage report:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

**Coverage by function:**
```bash
go tool cover -func=coverage.out
```

**Set coverage threshold:**
```bash
go test -cover ./... | grep "coverage:" | awk '{if($3+0 < 80) exit 1}'
```

## Continuous Integration

### GitHub Actions

Create `.github/workflows/test.yml`:
```yaml
name: Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -race -coverprofile=coverage.out ./...
    
    - name: Check coverage
      run: |
        go tool cover -func=coverage.out
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        echo "Coverage: $COVERAGE%"
        if (( $(echo "$COVERAGE < 80" | bc -l) )); then
          echo "Coverage is below 80%"
          exit 1
        fi
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

## Best Practices

### Test Organization

1. **One test file per source file** - `service.go` → `service_test.go`
2. **Group related tests** - Use subtests for variations
3. **Clear test names** - Describe what is being tested
4. **Arrange-Act-Assert** - Structure tests clearly

### Test Data

1. **Use test fixtures** - Create reusable test data
2. **Generate unique IDs** - Use `uuid.New()` for each test
3. **Avoid shared state** - Each test should be independent
4. **Use table-driven tests** - For testing multiple scenarios

### Mocking

1. **Mock external dependencies** - Databases, APIs, etc.
2. **Use interfaces** - Enable easy mocking
3. **Verify interactions** - Assert mock expectations
4. **Keep mocks simple** - Don't over-complicate

### Performance Testing

**Benchmark tests:**
```go
func BenchmarkAppointmentService_BookAppointment(b *testing.B) {
    repo := adapters.NewMemoryRepository()
    service := application.NewAppointmentService(repo, repo)
    
    serviceID := uuid.New()
    clientID := uuid.New()
    startTime := "2024-01-15T10:00:00Z"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.BookAppointment(context.Background(), serviceID, clientID, startTime)
    }
}
```

**Run benchmarks:**
```bash
go test -bench=. ./...
```