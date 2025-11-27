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

### Autenticación
- `POST /api/auth/login` — Login usuario (valida credenciales, retorna JWT)
    - Nota: si el usuario ya tiene una sesión activa (token vigente), el servidor responde 409 Conflict y no permitirá un nuevo login hasta cerrar sesión.
- `POST /api/auth/logout` — Logout (invalida el token y la sesión activa) **requiere JWT**

### Conexiones (stub)
- `GET /api/connections` — Stub, responde NotImplemented

### Auditorías (audits)
Rutas de auditoría ahora están agrupadas por gestor y siguen el patrón `/api/db/{gestor}/audits`.

- `POST /api/db/{gestor}/audits/execute` — Ejecuta una auditoría usando la conexión activa del usuario para `{gestor}` (ejecuta scripts de control seleccionados o por control). **requiere JWT**
- `GET /api/db/{gestor}/audits/:id` — Recupera el detalle de una auditoría y los resultados por script (audit run). **requiere JWT**

### Administración (admin)
- `GET /api/admin/sessions` — Lista usuarios con sesión activa y sus tokens. Requiere rol `admin` y token Bearer.

Ejemplo (usar token de admin en Authorization header):

Request:

```http
GET http://localhost:8000/api/admin/sessions
Authorization: Bearer <admin-jwt-token>
```

Respuesta esperada (200):

```json
{
    "sessions": [
        { "session_id": 10, "user_id": 1, "username": "admin", "email": "admin@example.com", "token": "eyJ...", "expires_at": "2025-11-26T20:00:00Z", "created_at": "2025-11-25T20:00:00Z" },
        { "session_id": 11, "user_id": 2, "username": "jane", "email": "jane@example.com", "token": "eyJ...", "expires_at": "2025-11-26T21:00:00Z", "created_at": "2025-11-25T21:00:00Z" }
    ]
}
```

Nota: almacenar tokens en claro puede no ser deseable en producción — considerar almacenar hash derivadas y dar a administradores sólo la visibilidad que realmente necesitan.

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

    3) Cerrar sesión (logout)

    ```http
    POST http://localhost:8000/api/auth/logout
    Authorization: Bearer <token>
    ```

    Respuesta esperada (200):

    ```json
    { "message": "logged out" }
    ```

- Paso 2 — Usar el token en Postman
    - Dentro de Postman, para cualquier endpoint protegido por JWT (por ejemplo, auditorías), en la pestaña `Headers` agregar:
        - Key: `Authorization`
        - Value: `Bearer <jwt-token>`
    - Alternativamente en Postman -> Authorization elegir tipo `Bearer Token` y pegar el token.

### Ejemplo: Ejecutar auditoría (POST /api/db/{gestor}/audits/execute)

 - URL: `POST http://localhost:8000/api/db/mssql/audits/execute`
- Headers:
    - `Content-Type: application/json`
    - `Authorization: Bearer <token>`
- Body (raw JSON) ejemplos:

1) Ejecutar por script IDs explícitos (usa la conexión activa para el gestor)

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

3) Ejecutar auditoría completa (ignora control_ids / script_ids)

```json
{
    "full_audit": true,
    "database": "master"
}
```

Notes:
- Manual controls (controls whose scripts are of type `manual`) are not executed on the server. For the purpose of the audit response, manual controls are considered "passed" (they require manual verification by the auditor) and are included in the `manual_count` field of the response.
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

### Ejemplo: Obtener auditoría (GET /api/db/{gestor}/audits/:id)

 - URL: `GET http://localhost:8000/api/db/mssql/audits/42` (reemplaza 42 por el id retornado en el POST)
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

### Conexiones (DB por gestor)

Estas rutas permiten a un usuario registrar una conexión activa a un servidor SQL para un gestor específico, listarlas y cerrarlas.

- `POST /api/db/{gestor}/open` — Crear/abrir una conexión activa para el gestor indicado **requiere JWT**
    - Body (JSON):
        ```json
        {
            "server": "dbserver.local",
            "port": "1433",
            "db_user": "dbuser",
            "password": "P@ssw0rd"
        }
        ```
        - Respuesta (200): contiene la información de la conexión (sin contraseñas en claro).
            La respuesta y la tabla `active_connections` ahora incluyen el campo `manager` que representa el gestor/driver lógico al que pertenece la conexión.
    - Seguridad: la contraseña se cifra antes de almacenarse en la base de datos (AES-GCM con la clave en `ENCRYPTION_KEY`).

- `GET /api/db/connections` — Obtener la lista de conexiones activas del usuario en todos los gestores (1 por gestor máximo) **requiere JWT**

- `GET /api/db/{gestor}/connection` — Obtener la conexión activa del usuario para el gestor indicado **requiere JWT**

- `DELETE /api/db/{gestor}/close` — Cerrar / eliminar la conexión activa del usuario para el gestor indicado **requiere JWT**

Ejemplo rápido en Postman para conectar (gestor mssql):

```http
POST http://localhost:8000/api/db/mssql/open
Authorization: Bearer <token>
Content-Type: application/json

{ "server": "localhost", "port": "1433", "db_user": "sa", "password": "YourStrong!Passw0rd" }
```

Respuesta esperada (200):

```json
{
    "connection": {
        "id": 2,
        "user_id": 5,
        "manager": "mssql",
        "driver": "mssql",
        "server": "localhost",
        "db_user": "sa",
        "is_connected": true,
        "last_connected": "2025-11-25T21:00:00Z"
    }
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