# API Documentation

## Base URL
`http://localhost:8000/api`

## Endpoints Actuales

### Health
- `GET /health` — Estado del backend (fuera de /api)
- `GET /` — Mensaje de bienvenida (fuera de /api)

### Swagger (stub)
- `GET /api/swagger` — Stub, no implementado

### Usuarios
- `GET /api/users/health` — Health de usuarios
- `POST /api/users/register` — Registrar usuario (stub, no guarda en BD)
- `POST /api/users/login` — Login usuario (stub, devuelve token fijo)

### Autenticación (legacy)
- `POST /api/auth/login` — Login usuario (stub, igual que /api/users/login)

### Conexiones (stub)
- `GET /api/connections` — Stub, responde NotImplemented

## Ejemplo de respuesta de stub

- `/api/users/register`:
```json
{
  "message": "user registered (stub)"
}
```
- `/api/users/login` o `/api/auth/login`:
```json
{
  "token": "stub-token"
}
```
- `/api/connections`:
```json
{
  "error": "connections endpoint not implemented yet"
}
```

## Notas
- Todos los endpoints reales de negocio (queries, conexiones, roles, etc.) aún no están implementados.
- Los endpoints actuales solo sirven para probar la estructura y evitar errores 404.
- Cuando se implementen los casos de uso reales, se actualizará esta documentación.

---

## Endpoints Planeados (según diseño original)

### Authentication

#### POST /auth/login
Login user and get JWT token.

Request:
```json
{
    "email": "user@example.com",
    "password": "password123"
}
```

Response:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
        "id": 1,
        "email": "user@example.com",
        "roles": ["user"]
    }
}
```

### Connections

#### POST /connections
Create new SQL Server connection.

Request:
```json
{
    "driver": "sqlserver",
    "server": "localhost",
    "port": "1433",
    "user": "sa",
    "password": "YourStrong!Passw0rd",
    "database": "master"
}
```

Response:
```json
{
    "id": 1,
    "isConnected": true,
    "server": "localhost",
    "database": "master"
}
```

#### GET /connections
List active connections.

Response:
```json
{
    "connections": [
        {
            "id": 1,
            "server": "localhost",
            "database": "master",
            "isConnected": true,
            "lastConnected": "2025-11-07T10:00:00Z"
        }
    ]
}
```

#### DELETE /connections/{id}
Close connection.

### Queries

#### POST /queries/execute
Execute SQL query.

Request:
```json
{
    "sql": "SELECT * FROM users",
    "database": "master",
    "pageSize": 100
}
```

Response:
```json
{
    "columns": ["id", "name", "email"],
    "types": ["int", "varchar", "varchar"],
    "rows": [
        [1, "John", "john@example.com"],
        [2, "Jane", "jane@example.com"]
    ],
    "hasMoreRows": false,
    "pageSize": 100,
    "page": 1
}
```

#### GET /queries/history
Get query execution history.

Query Parameters:
- page (default: 1)
- pageSize (default: 10)
- startDate (optional)
- endDate (optional)

Response:
```json
{
    "queries": [
        {
            "id": 1,
            "sql": "SELECT * FROM users",
            "status": "completed",
            "startTime": "2025-11-07T10:00:00Z",
            "endTime": "2025-11-07T10:00:01Z",
            "rowsAffected": 10
        }
    ],
    "total": 50,
    "page": 1,
    "pageSize": 10
}
```

### Users

#### POST /users
Create new user (admin only).

Request:
```json
{
    "email": "newuser@example.com",
    "password": "password123",
    "role": "user"
}
```

#### PUT /users/{id}
Update user (admin or self).

Request:
```json
{
    "email": "updated@example.com",
    "password": "newpassword123"
}
```

### Roles

#### POST /users/{id}/roles
Assign role to user (admin only).

Request:
```json
{
    "role": "manager"
}
```

#### GET /roles
List available roles.

Response:
```json
{
    "roles": [
        {
            "id": 1,
            "name": "admin",
            "permissions": [
                "manage_users",
                "execute_queries",
                "manage_connections"
            ]
        }
    ]
}
```

## Error Responses

All errors follow this format:
```json
{
    "error": "Error description"
}
```

Common HTTP Status Codes:
- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 429: Too Many Requests
- 500: Internal Server Error

## Rate Limiting

- 100 requests per minute per user
- 5 concurrent queries per user
- Headers included in response:
  - X-RateLimit-Limit
  - X-RateLimit-Remaining
  - X-RateLimit-Reset

## Pagination

Endpoints that return lists support pagination:
- page: Page number (starts at 1)
- pageSize: Items per page
- Response includes total count and pagination info

## Data Types

When executing queries, the following type mappings are used:

SQL Server -> JSON:
- int -> number
- varchar/nvarchar -> string
- datetime -> string (ISO 8601)
- bit -> boolean
- decimal/numeric -> string
- binary/varbinary -> base64 string

## Security

- All requests must use HTTPS in production
- Passwords must be at least 8 characters
- SQL injection protection enabled
- Rate limiting per user/IP
- Token expiration: 24 hours