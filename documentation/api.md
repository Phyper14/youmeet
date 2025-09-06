# API Documentation

## Overview

YouMeet provides a RESTful API for managing appointments and services. All endpoints accept and return JSON data.

## Base URL

```
http://localhost:8080
```

## Endpoints

### Book Appointment

Books a new appointment for a client.

**Endpoint:** `POST /appointments`

**Request Body:**
```json
{
  "service_id": "uuid",
  "client_id": "uuid", 
  "start_time": "2024-01-15T10:00:00Z"
}
```

**Response:**
```json
{
  "id": "uuid",
  "service_id": "uuid",
  "client_id": "uuid",
  "start_time": "2024-01-15T10:00:00Z",
  "status": "scheduled"
}
```

**Status Codes:**
- `200 OK` - Appointment booked successfully
- `400 Bad Request` - Invalid request body or parameters
- `500 Internal Server Error` - Server error

**Example:**
```bash
curl -X POST http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "service_id": "123e4567-e89b-12d3-a456-426614174000",
    "client_id": "987fcdeb-51a2-43d1-9c4b-123456789abc",
    "start_time": "2024-01-15T10:00:00Z"
  }'
```

### Get Appointments

Retrieves all appointments for a specific client.

**Endpoint:** `GET /appointments/{clientId}`

**Path Parameters:**
- `clientId` (string, required) - UUID of the client

**Response:**
```json
[
  {
    "id": "uuid",
    "service_id": "uuid",
    "client_id": "uuid",
    "start_time": "2024-01-15T10:00:00Z",
    "status": "scheduled"
  }
]
```

**Status Codes:**
- `200 OK` - Appointments retrieved successfully
- `400 Bad Request` - Invalid client ID format
- `500 Internal Server Error` - Server error

**Example:**
```bash
curl http://localhost:8080/appointments/987fcdeb-51a2-43d1-9c4b-123456789abc
```

## Data Models

### Appointment

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique appointment identifier |
| service_id | UUID | ID of the booked service |
| client_id | UUID | ID of the client |
| start_time | DateTime | Appointment start time (RFC3339 format) |
| status | String | Appointment status (e.g., "scheduled") |

### Service

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique service identifier |
| name | String | Service name |
| duration | Duration | Service duration |
| price | Float64 | Service price |
| provider_id | UUID | ID of the service provider |

### User

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique user identifier |
| name | String | User's full name |
| email | String | User's email address |

## Error Handling

All errors return a JSON response with an error message:

```json
{
  "error": "Error description"
}
```

Common error scenarios:
- Invalid UUID format in path parameters
- Missing required fields in request body
- Invalid date format for start_time (must be RFC3339)
- Internal server errors

## Date Format

All datetime fields use RFC3339 format: `2006-01-02T15:04:05Z07:00`

Examples:
- `2024-01-15T10:00:00Z` (UTC)
- `2024-01-15T10:00:00-05:00` (EST)