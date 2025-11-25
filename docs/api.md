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
- `POST /api/users/register` — Registrar usuario (crea usuario en BD y retorna JWT)
- `POST /api/users/login` — Login usuario (valida credenciales, retorna JWT)

### Autenticación
- `POST /api/auth/login` — Login usuario (valida credenciales, retorna JWT)

### Conexiones (stub)
- `GET /api/connections` — Stub, responde NotImplemented

### Auditorías (audits)
- `POST /api/audits/execute` — Ejecuta una auditoría (ejecuta scripts de control seleccionados o por control) **requiere JWT**
- `GET /api/audits/:id` — Recupera el detalle de una auditoría y los resultados por script **requiere JWT**

## Ejemplo de respuesta de stub

`/api/users/register` (success):
```json
{
    "token": "<jwt-token>",
    "user": {
        "id": 1,
        "username": "jdoe",
        "email": "jdoe@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "role": "cliente"
    }
}
```
`/api/users/login` or `/api/auth/login` (success):
```json
{
    "token": "<jwt-token>",
    "user": { /* same shape as above */ }
}
```

### Cómo usar Postman (rápido)

- Paso 1 — Obtener token JWT (register/login)
    1. Abrir Postman
    2. Crear request POST `http://localhost:8000/api/users/register` o `POST http://localhost:8000/api/auth/login`
    3. En el body usar `raw` JSON con contenido similar a:

```json
{
    "username": "tester",
    "email": "tester@example.com",
    "password": "password123"
}
```

    4. En la respuesta, copiar el campo `token` (JWT).

- Paso 2 — Usar el token en Postman
    - Dentro de Postman, para cualquier endpoint protegido por JWT (por ejemplo, auditorías), en la pestaña `Headers` agregar:
        - Key: `Authorization`
        - Value: `Bearer <jwt-token>`
    - Alternativamente en Postman -> Authorization elegir tipo `Bearer Token` y pegar el token.

### Ejemplo: Ejecutar auditoría (POST /api/audits/execute)

- URL: `POST http://localhost:8000/api/audits/execute`
- Headers:
    - `Content-Type: application/json`
    - `Authorization: Bearer <token>`
- Body (raw JSON) ejemplos:

1) Ejecutar por script IDs explícitos

```json
{
    "script_ids": [1, 2, 10],
    "database": "master"
}
```

2) Ejecutar por control IDs (agrupa scripts)

```json
{
    "control_ids": [3, 4],
    "database": "master"
}
```

Respuesta esperada (200):

```json
{
    "total": 3,
    "passed": 3,
    "failed": 0,
    "scripts": [
        {"script_id": 1, "control_id": 3, "control_type": "simple", "query_sql": "SELECT 1", "passed": true},
        {"script_id": 2, "control_id": 4, "control_type": "index", "query_sql": "SELECT COUNT(*) FROM users", "passed": true}
    ],
    "audit_run_id": 42
}
```

Nota: la respuesta contiene `audit_run_id` cuando la persistencia de auditorías está habilitada; puedes usar este id para consultar el run detallado.

### Ejemplo: Obtener auditoría (GET /api/audits/:id)

- URL: `GET http://localhost:8000/api/audits/42` (reemplaza 42 por el id retornado en el POST)
- Headers: `Authorization: Bearer <token>`

Respuesta esperada (200):

```json
{
    "audit": {
        "id": 42,
        "user_id": 1,
        "mode": "partial",
        "database": "master",
        "total": 3,
        "passed": 3,
        "failed": 0,
        "status": "completed",
        "controls": "[3,4]",
        "started_at": "2025-11-25T14:00:00Z",
        "finished_at": "2025-11-25T14:00:01Z"
    },
    "result": {
        "total": 3,
        "passed": 3,
        "failed": 0,
        "scripts": [ /* per-script results */ ]
    }
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