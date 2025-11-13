# JWT Authentication - Testing Guide

Esta guía proporciona ejemplos de cómo probar los endpoints de autenticación JWT implementados.

---

## Información de Configuración

- **Base URL**: `http://localhost:8000` (o el puerto configurado en SERVER_PORT)
- **JWT Secret**: Configurado via `JWT_SECRET` env var (default: "change-me-in-production")
- **Token Expiry**: 24 horas
- **Database**: MySQL en `localhost:3306` (configurable via MYSQL_* env vars)

---

## Prerequisites

Asegúrate de que:
1. El servidor está corriendo: `go run ./cmd/server/main.go`
2. MySQL está disponible con la base de datos `microsql_ago` creada
3. Las migraciones han corrido (tablas de usuarios existen)
4. Hay al menos un usuario registrado en la DB (o crear uno manualmente)

---

## Test Cases

### 1. **Login - Obtener JWT Token**

**Endpoint**: `POST /api/auth/login`

**Request**:
```bash
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**Expected Response** (201 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzMxNDQzNzk5LCJpYXQiOjE3MzEzNTczOTksIm5iZiI6MTczMTM1NzM5OX0.signature...",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "first_name": "Admin",
    "last_name": "User",
    "role": "admin"
  }
}
```

**Error Responses**:
- `400 Bad Request`: JSON malformado o campos faltantes
- `401 Unauthorized`: Credenciales inválidas
  ```json
  { "error": "invalid credentials" }
  ```

---

### 2. **Access Protected Endpoint - Execute Audit (con token)**

**Endpoint**: `POST /api/audits/execute`

**Request** (con token válido):
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."

curl -X POST http://localhost:8000/api/audits/execute \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "control_ids": [1, 2],
    "database": "AuditedDatabase"
  }'
```

**Expected Response** (200 OK):
```json
{
  "total": 2,
  "passed": 1,
  "failed": 1,
  "scripts": [
    {
      "script_id": 1,
      "control_id": 1,
      "control_type": "S",
      "query_sql": "SELECT ... ",
      "passed": true
    },
    {
      "script_id": 2,
      "control_id": 2,
      "control_type": "S",
      "query_sql": "SELECT ... ",
      "passed": false,
      "error": "timeout or error message"
    }
  ]
}
```

---

### 3. **Access Protected Endpoint - Sin Token**

**Request** (sin Authorization header):
```bash
curl -X POST http://localhost:8000/api/audits/execute \
  -H "Content-Type: application/json" \
  -d '{
    "control_ids": [1, 2],
    "database": "AuditedDatabase"
  }'
```

**Expected Response** (401 Unauthorized):
```json
{
  "error": "no token provided"
}
```

---

### 4. **Access Protected Endpoint - Token Inválido**

**Request** (con token malformado):
```bash
curl -X POST http://localhost:8000/api/audits/execute \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid.token.here" \
  -d '{
    "control_ids": [1, 2],
    "database": "AuditedDatabase"
  }'
```

**Expected Response** (401 Unauthorized):
```json
{
  "error": "invalid or expired token: failed to parse token: ..."
}
```

---

### 5. **Access Protected Endpoint - Token Expirado**

Una vez que un token expire (después de 24 horas), intentar usarlo resultará en:

**Expected Response** (401 Unauthorized):
```json
{
  "error": "invalid or expired token: token is not valid"
}
```

---

### 6. **Bad Authorization Header Format**

**Request** (sin "Bearer" prefix):
```bash
curl -X POST http://localhost:8000/api/audits/execute \
  -H "Content-Type: application/json" \
  -H "Authorization: $TOKEN" \
  -d '{
    "control_ids": [1],
    "database": "DB"
  }'
```

**Expected Response** (401 Unauthorized):
```json
{
  "error": "invalid authorization header"
}
```

---

## Test Flow Completo (paso a paso)

```bash
#!/bin/bash

# 1. Login y obtener token
echo "=== LOGIN ==="
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

echo "Login Response:"
echo "$LOGIN_RESPONSE" | jq .

# Extraer token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')
echo "Token: $TOKEN"

# 2. Usar token para acceder a endpoint protegido
echo -e "\n=== EXECUTE AUDIT (with valid token) ==="
curl -s -X POST http://localhost:8000/api/audits/execute \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "control_ids": [1],
    "database": "TestDB"
  }' | jq .

# 3. Intentar sin token
echo -e "\n=== EXECUTE AUDIT (without token) ==="
curl -s -X POST http://localhost:8000/api/audits/execute \
  -H "Content-Type: application/json" \
  -d '{
    "control_ids": [1],
    "database": "TestDB"
  }' | jq .

# 4. Intentar con token inválido
echo -e "\n=== EXECUTE AUDIT (with invalid token) ==="
curl -s -X POST http://localhost:8000/api/audits/execute \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid" \
  -d '{
    "control_ids": [1],
    "database": "TestDB"
  }' | jq .
```

---

## JWT Token Anatomy (Decode)

Para verificar el contenido de un token (sin verificar firma):

```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzMxNDQzNzk5LCJpYXQiOjE3MzEzNTczOTksIm5iZiI6MTczMTM1NzM5OX0.signature"

# Usa jq para decodificar (nota: no verifica firma)
echo "$TOKEN" | cut -d. -f2 | base64 -D | jq .
```

**Ejemplo de claims decodificadas**:
```json
{
  "user_id": 1,
  "username": "admin",
  "role": "admin",
  "exp": 1731443799,    // Unix timestamp de expiración
  "iat": 1731357399,    // Unix timestamp de emisión
  "nbf": 1731357399     // Unix timestamp notBefore
}
```

---

## Postman Collection

Puedes importar esta colección en Postman para probar la API:

```json
{
  "info": {
    "name": "MicroSQL-AGo JWT Auth",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Login",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\"username\": \"admin\", \"password\": \"admin123\"}"
        },
        "url": {
          "raw": "http://localhost:8000/api/auth/login",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8000",
          "path": ["api", "auth", "login"]
        }
      }
    },
    {
      "name": "Execute Audit (Protected)",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          },
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\"control_ids\": [1, 2], \"database\": \"AuditedDB\"}"
        },
        "url": {
          "raw": "http://localhost:8000/api/audits/execute",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8000",
          "path": ["api", "audits", "execute"]
        }
      }
    }
  ]
}
```

---

## Troubleshooting

### "no token provided"
- El header `Authorization` está vacío o no se envió
- Verifica que usas: `Authorization: Bearer <token>`

### "invalid or expired token"
- El token es malformado
- El token ha expirado (pasaron más de 24 horas)
- El JWT_SECRET en el servidor no coincide con el que generó el token

### "invalid credentials"
- Username no existe
- Password es incorrecta
- Usuario está inactivo (`is_active = false`)

### 500 Internal Server Error en login
- Base de datos no está disponible
- Tabla de usuarios no existe
- Error al hashear/validar password

---

## Notas de Seguridad

- **NEVER** envíes JWT en URL (siempre en header Authorization o en body POST si es absolutamente necesario)
- Cambia `JWT_SECRET` en producción a un valor fuerte y único
- Usa HTTPS en producción para evitar interception de tokens
- Implementa rate limiting en endpoints de login para prevenir brute force
- Considera agregar refresh tokens para sesiones largas (fuera de scope actual)
