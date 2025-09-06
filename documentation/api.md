# API Documentation - YouMeet

## Visão Geral

A API YouMeet é uma RESTful API construída em Go que permite o gerenciamento completo de agendamentos para prestadores de serviços.

## Base URL

```
http://localhost:8080
```

## Autenticação

### POST /auth/register

Registra um novo usuário no sistema.

**Request Body:**
```json
{
  "name": "João Silva",
  "email": "joao@email.com",
  "password": "senha123",
  "role": "client"
}
```

**Roles disponíveis:**
- `client` - Cliente que agenda serviços
- `company` - Empresa que oferece serviços
- `professional` - Profissional autônomo

**Response (201):**
```json
{
  "message": "Usuário criado com sucesso",
  "user": {
    "id": "uuid",
    "name": "João Silva",
    "email": "joao@email.com",
    "role": "client"
  }
}
```

### POST /auth/login

Autentica um usuário existente.

**Request Body:**
```json
{
  "email": "joao@email.com",
  "password": "senha123"
}
```

**Response (200):**
```json
{
  "token": "auth-token-uuid",
  "user": {
    "id": "uuid",
    "name": "João Silva",
    "email": "joao@email.com",
    "role": "client"
  }
}
```

## Agendamentos

### POST /appointments

Cria um novo agendamento.

**Request Body:**
```json
{
  "client_id": "client-uuid",
  "professional_id": "professional-uuid",
  "service_id": "service-uuid",
  "start_time": "2024-01-15T10:00:00Z",
  "end_time": "2024-01-15T11:00:00Z"
}
```

**Response (201):**
```json
{
  "message": "Agendamento criado com sucesso",
  "appointment": {
    "id": "appointment-uuid",
    "client_id": "client-uuid",
    "professional_id": "professional-uuid",
    "service_id": "service-uuid",
    "start_time": "2024-01-15T10:00:00Z",
    "end_time": "2024-01-15T11:00:00Z",
    "status": "scheduled"
  }
}
```

### GET /appointments/{clientId}

Lista todos os agendamentos de um cliente.

**Response (200):**
```json
{
  "appointments": [
    {
      "id": "appointment-uuid",
      "client_id": "client-uuid",
      "professional_id": "professional-uuid",
      "service_id": "service-uuid",
      "start_time": "2024-01-15T10:00:00Z",
      "end_time": "2024-01-15T11:00:00Z",
      "status": "scheduled"
    }
  ]
}
```

## Códigos de Status

- `200` - OK
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

## Exemplos de Uso

### Registrar um cliente
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Maria Santos",
    "email": "maria@email.com",
    "password": "senha123",
    "role": "client"
  }'
```

### Fazer login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@email.com",
    "password": "senha123"
  }'
```

### Criar agendamento
```bash
curl -X POST http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "client-uuid",
    "professional_id": "professional-uuid",
    "service_id": "service-uuid",
    "start_time": "2024-01-15T10:00:00Z",
    "end_time": "2024-01-15T11:00:00Z"
  }'
```