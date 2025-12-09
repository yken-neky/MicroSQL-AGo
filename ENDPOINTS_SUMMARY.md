# Resumen de Endpoints - MicroSQL AGo

Este documento contiene un resumen completo de todas las funcionalidades implementadas en los endpoints del proyecto MicroSQL AGo.

---

## 📋 Tabla de Contenidos

1. [Endpoints Generales](#endpoints-generales)
2. [Endpoints de Autenticación](#endpoints-de-autenticación)
3. [Endpoints de Usuarios](#endpoints-de-usuarios)
4. [Endpoints de Conexiones a Bases de Datos](#endpoints-de-conexiones-a-bases-de-datos)
5. [Endpoints de Auditorías](#endpoints-de-auditorías)
6. [Endpoints de Administración](#endpoints-de-administración)

---

## 🔧 Endpoints Generales

### `GET /`
**Descripción:** Endpoint raíz que devuelve información básica del servicio.

**Autenticación:** No requiere autenticación

**Respuesta:**
```json
{
  "service": "MicroSQL AGo backend",
  "status": "ok"
}
```

---

### `GET /health`
**Descripción:** Endpoint de health check para verificar el estado del servicio.

**Autenticación:** No requiere autenticación

**Respuesta:**
```json
{
  "status": "ok"
}
```

---

### `GET /api/swagger`
**Descripción:** Endpoint para información de Swagger (actualmente no generado).

**Autenticación:** No requiere autenticación

**Respuesta:**
```json
{
  "swagger": "not generated"
}
```

---

### `GET /api/users/health`
**Descripción:** Health check específico para el módulo de usuarios.

**Autenticación:** No requiere autenticación

**Respuesta:**
```json
{
  "ok": true
}
```

---

## 🔐 Endpoints de Autenticación

### `POST /api/auth/login`
**Descripción:** Inicia sesión de un usuario y genera un token JWT. Implementa política de sesión única (un usuario solo puede tener una sesión activa a la vez).

**Autenticación:** No requiere autenticación

**Request Body:**
```json
{
  "username": "string (requerido)",
  "password": "string (requerido)"
}
```

**Validaciones:**
- Verifica credenciales (username y password)
- Verifica que el usuario esté activo (`is_active = true`)
- Verifica que no exista una sesión activa previa (si existe y no está expirada, rechaza el login)
- Si existe una sesión expirada, la marca como inactiva y permite el nuevo login

**Respuesta Exitosa (200):**
```json
{
  "token": "jwt_token_string",
  "user": {
    "id": 1,
    "username": "usuario",
    "email": "usuario@example.com",
    "first_name": "Nombre",
    "last_name": "Apellido",
    "role": "user"
  }
}
```

**Errores:**
- `400`: Error en el formato del request
- `401`: Credenciales inválidas o usuario inactivo
- `409`: Usuario ya tiene una sesión activa
- `500`: Error interno del servidor

**Funcionalidades adicionales:**
- Actualiza `last_login` del usuario al hacer login exitoso
- Crea una sesión en la base de datos con expiración de 24 horas
- Genera token JWT con información del usuario (ID, username, role)

---

### `POST /api/auth/logout`
**Descripción:** Cierra la sesión del usuario actual invalidando el token JWT presentado.

**Autenticación:** Requiere token JWT válido (Bearer token)

**Headers:**
```
Authorization: Bearer <token>
```

**Respuesta Exitosa (200):**
```json
{
  "message": "logged out"
}
```

**Funcionalidades:**
- Marca la sesión como inactiva (`is_active = false`)
- Es idempotente: si no hay sesión activa, devuelve éxito de todas formas

---

## 👤 Endpoints de Usuarios

### `POST /api/users/register`
**Descripción:** Registra un nuevo usuario en el sistema. Crea el usuario con rol "user" por defecto y genera un token JWT automáticamente.

**Autenticación:** No requiere autenticación

**Request Body:**
```json
{
  "username": "string (requerido, min 3, max 150 caracteres)",
  "first_name": "string (requerido, min 3, max 150 caracteres)",
  "last_name": "string (requerido, min 3, max 150 caracteres)",
  "email": "string (requerido, formato email válido)",
  "password": "string (requerido, mínimo 8 caracteres)"
}
```

**Validaciones:**
- Verifica que el username no esté tomado
- Verifica que el email no esté registrado
- Valida formato de email
- Valida longitud mínima de password (8 caracteres)

**Respuesta Exitosa (201):**
```json
{
  "token": "jwt_token_string",
  "user": {
    "id": 1,
    "username": "nuevo_usuario",
    "email": "nuevo@example.com",
    "first_name": "Nombre",
    "last_name": "Apellido",
    "role": "user"
  }
}
```

**Funcionalidades:**
- Hashea la contraseña con bcrypt antes de almacenarla
- Crea el usuario con `is_active = true` y `role = "user"`
- Asigna automáticamente el rol "user" (role_id = 3) en la tabla `user_roles`
- Genera token JWT automáticamente para el nuevo usuario
- Actualiza `last_login` al momento de registro

**Errores:**
- `400`: Error en validación de campos
- `409`: Username o email ya existe
- `500`: Error interno del servidor

---

## 🗄️ Endpoints de Conexiones a Bases de Datos

Todos los endpoints de conexión requieren autenticación JWT y están bajo el prefijo `/api/db`.

### `GET /api/db/connections`
**Descripción:** Lista todas las conexiones activas del usuario autenticado, sin importar el gestor de base de datos.

**Autenticación:** Requiere token JWT válido

**Respuesta Exitosa (200):**
```json
{
  "connections": [
    {
      "id": 1,
      "user_id": 1,
      "manager": "mssql",
      "driver": "mssql",
      "server": "localhost",
      "db_user": "sa",
      "is_connected": true,
      "last_connected": "2024-01-01T10:00:00Z",
      "last_disconnected": null
    }
  ]
}
```

---

### `POST /api/db/:manager/open`
**Descripción:** Abre una nueva conexión a un servidor de base de datos usando el gestor especificado en la URL.

**Autenticación:** Requiere token JWT válido

**Parámetros de URL:**
- `manager`: Gestor de base de datos (`pgsql`, `oracle`, `mysql`, `mssql`, `otro`)

**Request Body:**
```json
{
  "manager": "string (requerido)",
  "driver": "string (requerido, ej: 'mssql')",
  "server": "string (requerido, host o IP)",
  "port": "string (requerido, ej: '1433')",
  "db_user": "string (requerido)",
  "password": "string (requerido, texto plano sobre TLS)"
}
```

**Validaciones:**
- Valida que el manager sea uno de los soportados
- Intenta establecer conexión real con el servidor de base de datos
- Encripta la contraseña antes de almacenarla

**Respuesta Exitosa (200):**
```json
{
  "connection": {
    "id": 1,
    "user_id": 1,
    "manager": "mssql",
    "driver": "mssql",
    "server": "localhost",
    "db_user": "sa",
    "is_connected": true,
    "last_connected": "2024-01-01T10:00:00Z"
  }
}
```

**Funcionalidades:**
- Crea una conexión activa en la base de datos
- Encripta y almacena las credenciales de forma segura
- Verifica la conexión antes de persistirla
- Registra el timestamp de conexión

**Errores:**
- `400`: Manager no soportado, datos inválidos, o error de conexión
- `500`: Error interno del servidor

---

### `DELETE /api/db/:manager/close`
**Descripción:** Cierra la conexión activa del usuario para el gestor especificado.

**Autenticación:** Requiere token JWT válido

**Parámetros de URL:**
- `manager`: Gestor de base de datos (`pgsql`, `oracle`, `mysql`, `mssql`, `otro`)

**Respuesta Exitosa (200):**
```json
{
  "message": "disconnected"
}
```

**Funcionalidades:**
- Cierra la conexión física con el servidor de base de datos
- Marca la conexión como desconectada (`is_connected = false`)
- Registra el timestamp de desconexión

**Errores:**
- `400`: Manager no soportado o no hay conexión activa
- `500`: Error interno del servidor

---

### `GET /api/db/:manager/connection`
**Descripción:** Obtiene la información de la conexión activa del usuario para el gestor especificado.

**Autenticación:** Requiere token JWT válido

**Parámetros de URL:**
- `manager`: Gestor de base de datos (`pgsql`, `oracle`, `mysql`, `mssql`, `otro`)

**Respuesta Exitosa (200):**
```json
{
  "connection": {
    "id": 1,
    "user_id": 1,
    "driver": "mssql",
    "server": "localhost",
    "db_user": "sa",
    "is_connected": true,
    "last_connected": "2024-01-01T10:00:00Z",
    "last_disconnected": null
  }
}
```

**Errores:**
- `404`: No hay conexión activa para ese gestor
- `400`: Manager no soportado
- `500`: Error interno del servidor

---

## 🔍 Endpoints de Auditorías

Los endpoints de auditorías están bajo `/api/db/:manager/audits` y requieren autenticación JWT.

### `POST /api/db/:manager/audits/execute`
**Descripción:** Ejecuta una auditoría parcial o completa sobre la base de datos conectada. Permite ejecutar controles de auditoría definidos en el sistema.

**Autenticación:** Requiere token JWT válido

**Parámetros de URL:**
- `manager`: Gestor de base de datos (`pgsql`, `oracle`, `mysql`, `mssql`, `otro`)

**Request Body:**
```json
{
  "control_ids": [1, 2, 3],  // IDs de controles específicos (opcional)
  "execute_all": false       // Si es true, ejecuta todos los controles
}
```

**Respuesta Exitosa (200):**
```json
{
  "audit_run_id": 123,
  "status": "completed",
  "results": [
    {
      "control_id": 1,
      "control_name": "Verificar usuarios sin contraseña",
      "status": "passed",
      "details": "..."
    }
  ]
}
```

**Funcionalidades:**
- Ejecuta scripts SQL de auditoría sobre la base de datos conectada
- Registra los resultados de cada control ejecutado
- Crea un registro de ejecución de auditoría (`audit_run`)
- Almacena los resultados detallados de cada script ejecutado
- Valida que el usuario tenga una conexión activa para el gestor especificado

**Errores:**
- `400`: Request inválido o sin conexión activa
- `500`: Error al ejecutar la auditoría

---

### `GET /api/db/:manager/audits/:id`
**Descripción:** Obtiene los detalles de una ejecución de auditoría específica, incluyendo el estado y los resultados de cada control ejecutado.

**Autenticación:** Requiere token JWT válido

**Parámetros de URL:**
- `manager`: Gestor de base de datos
- `id`: ID de la ejecución de auditoría

**Respuesta Exitosa (200):**
```json
{
  "audit": {
    "id": 123,
    "user_id": 1,
    "manager": "mssql",
    "status": "completed",
    "started_at": "2024-01-01T10:00:00Z",
    "finished_at": "2024-01-01T10:05:00Z"
  },
  "result": {
    "results": [
      {
        "control_id": 1,
        "control_name": "Verificar usuarios sin contraseña",
        "status": "passed",
        "script_output": "...",
        "execution_time_ms": 150
      }
    ]
  }
}
```

**Validaciones:**
- Verifica que la auditoría pertenezca al usuario autenticado
- Solo permite acceso a auditorías del usuario que las ejecutó

**Errores:**
- `400`: ID de auditoría inválido
- `403`: La auditoría no pertenece al usuario
- `500`: Error interno del servidor

---

## 👨‍💼 Endpoints de Administración

Todos los endpoints de administración requieren:
- Autenticación JWT válida
- Rol de "admin"

Están bajo el prefijo `/api/admin`.

### `GET /api/admin/sessions`
**Descripción:** Lista todas las sesiones activas en el sistema con información de los usuarios asociados.

**Autenticación:** Requiere rol "admin"

**Respuesta Exitosa (200):**
```json
{
  "sessions": [
    {
      "session_id": 1,
      "user_id": 1,
      "username": "usuario",
      "email": "usuario@example.com",
      "token": "jwt_token_string",
      "expires_at": "2024-01-02T10:00:00Z",
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

---

### Gestión de Roles

#### `GET /api/admin/roles`
**Descripción:** Lista todos los roles disponibles en el sistema con sus permisos asociados.

**Autenticación:** Requiere rol "admin"

**Respuesta Exitosa (200):**
```json
{
  "roles": [
    {
      "id": 1,
      "name": "admin",
      "description": "Administrador del sistema",
      "permissions": [...]
    }
  ]
}
```

---

#### `POST /api/admin/roles`
**Descripción:** Crea un nuevo rol en el sistema.

**Autenticación:** Requiere rol "admin"

**Request Body:**
```json
{
  "name": "string (requerido)",
  "description": "string (opcional)"
}
```

**Respuesta Exitosa (201):**
```json
{
  "id": 4,
  "name": "auditor",
  "description": "Rol para usuarios auditores"
}
```

**Funcionalidades:**
- Registra la acción en el log de auditoría RBAC

---

#### `PUT /api/admin/roles/:id`
**Descripción:** Actualiza los metadatos de un rol existente.

**Autenticación:** Requiere rol "admin"

**Parámetros de URL:**
- `id`: ID del rol a actualizar

**Request Body:**
```json
{
  "name": "string (opcional)",
  "description": "string (opcional)"
}
```

**Respuesta Exitosa (200):** Rol actualizado

**Funcionalidades:**
- Registra la acción en el log de auditoría RBAC

---

#### `DELETE /api/admin/roles/:id`
**Descripción:** Elimina un rol del sistema.

**Autenticación:** Requiere rol "admin"

**Parámetros de URL:**
- `id`: ID del rol a eliminar

**Respuesta Exitosa (204):** Sin contenido

**Funcionalidades:**
- Registra la acción en el log de auditoría RBAC

---

### Gestión de Usuarios y Roles

#### `GET /api/admin/users`
**Descripción:** Lista todos los usuarios del sistema con sus roles asociados.

**Autenticación:** Requiere rol "admin"

**Respuesta Exitosa (200):**
```json
{
  "users": [
    {
      "id": 1,
      "username": "usuario",
      "email": "usuario@example.com",
      "is_active": true,
      "roles": ["user", "auditor"]
    }
  ]
}
```

---

#### `POST /api/admin/users/:id/roles`
**Descripción:** Asigna un rol a un usuario.

**Autenticación:** Requiere rol "admin"

**Parámetros de URL:**
- `id`: ID del usuario

**Request Body:**
```json
{
  "role_id": 2
}
```

**Respuesta Exitosa (200):**
```json
{
  "ok": true
}
```

**Funcionalidades:**
- Registra la acción en el log de auditoría RBAC

---

#### `DELETE /api/admin/users/:id/roles`
**Descripción:** Revoca un rol de un usuario.

**Autenticación:** Requiere rol "admin"

**Parámetros de URL:**
- `id`: ID del usuario

**Request Body:**
```json
{
  "role_id": 2
}
```

**Respuesta Exitosa (200):**
```json
{
  "ok": true
}
```

**Funcionalidades:**
- Registra la acción en el log de auditoría RBAC

---

### Gestión de Permisos

#### `GET /api/admin/permissions`
**Descripción:** Lista todos los permisos disponibles en el sistema.

**Autenticación:** Requiere rol "admin"

**Respuesta Exitosa (200):**
```json
{
  "permissions": [
    {
      "id": 1,
      "name": "execute_audit",
      "resource": "audits",
      "action": "execute",
      "description": "Permite ejecutar auditorías"
    }
  ]
}
```

---

#### `POST /api/admin/permissions`
**Descripción:** Crea un nuevo permiso.

**Autenticación:** Requiere rol "admin"

**Request Body:**
```json
{
  "name": "string (requerido)",
  "resource": "string (requerido)",
  "action": "string (requerido)",
  "description": "string (opcional)"
}
```

**Respuesta Exitosa (201):** Permiso creado

**Funcionalidades:**
- Registra la acción en el log de auditoría RBAC

---

#### `PUT /api/admin/permissions/:id`
**Descripción:** Actualiza un permiso existente.

**Autenticación:** Requiere rol "admin"

**Parámetros de URL:**
- `id`: ID del permiso

**Request Body:**
```json
{
  "name": "string (opcional)",
  "resource": "string (opcional)",
  "action": "string (opcional)",
  "description": "string (opcional)"
}
```

**Respuesta Exitosa (200):** Permiso actualizado

**Funcionalidades:**
- Registra la acción en el log de auditoría RBAC

---

#### `DELETE /api/admin/permissions/:id`
**Descripción:** Elimina un permiso del sistema.

**Autenticación:** Requiere rol "admin"

**Parámetros de URL:**
- `id`: ID del permiso

**Respuesta Exitosa (204):** Sin contenido

**Funcionalidades:**
- Registra la acción en el log de auditoría RBAC

---

### Asignación de Permisos a Roles

#### `POST /api/admin/roles/:id/permissions`
**Descripción:** Asigna un permiso a un rol.

**Autenticación:** Requiere rol "admin"

**Parámetros de URL:**
- `id`: ID del rol

**Request Body:**
```json
{
  "permission_id": 1
}
```

**Respuesta Exitosa (200):**
```json
{
  "ok": true
}
```

**Funcionalidades:**
- Registra la acción en el log de auditoría RBAC

---

#### `DELETE /api/admin/roles/:id/permissions`
**Descripción:** Revoca un permiso de un rol.

**Autenticación:** Requiere rol "admin"

**Parámetros de URL:**
- `id`: ID del rol

**Request Body:**
```json
{
  "permission_id": 1
}
```

**Respuesta Exitosa (200):**
```json
{
  "ok": true
}
```

**Funcionalidades:**
- Registra la acción en el log de auditoría RBAC

---

### Auditoría RBAC

#### `GET /api/admin/audit/rbac`
**Descripción:** Lista los logs de auditoría de acciones RBAC (creación/actualización/eliminación de roles, permisos, asignaciones, etc.).

**Autenticación:** Requiere rol "admin"

**Query Parameters:**
- `actor_id` (opcional): Filtrar por ID del actor
- `target_type` (opcional): Filtrar por tipo de objetivo (ej: "role", "permission", "user_role")
- `action` (opcional): Filtrar por acción (ej: "role.create", "permission.assign")
- `limit` (opcional, default: 100): Límite de resultados
- `offset` (opcional, default: 0): Offset para paginación

**Respuesta Exitosa (200):**
```json
{
  "logs": [
    {
      "id": 1,
      "actor_id": 1,
      "actor_name": "admin",
      "action": "role.create",
      "target_type": "role",
      "target_id": 4,
      "target_name": "auditor",
      "details": "Rol para usuarios auditores",
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

---

### Métricas del Sistema

#### `GET /api/admin/metrics/users`
**Descripción:** Obtiene métricas sobre los usuarios del sistema.

**Autenticación:** Requiere rol "admin"

**Respuesta Exitosa (200):**
```json
{
  "total_users": 100,
  "active_users": 85,
  "roles_distribution": [
    {
      "role": "user",
      "count": 80
    },
    {
      "role": "admin",
      "count": 5
    }
  ]
}
```

---

#### `GET /api/admin/metrics/connections`
**Descripción:** Obtiene métricas sobre las conexiones a bases de datos.

**Autenticación:** Requiere rol "admin"

**Respuesta Exitosa (200):**
```json
{
  "total_active_connections": 50,
  "currently_connected": 30,
  "total_connection_logs": 500
}
```

---

#### `GET /api/admin/metrics/audits`
**Descripción:** Obtiene métricas sobre las ejecuciones de auditorías.

**Autenticación:** Requiere rol "admin"

**Respuesta Exitosa (200):**
```json
{
  "total_runs": 200,
  "status_distribution": [
    {
      "status": "completed",
      "count": 180
    },
    {
      "status": "failed",
      "count": 20
    }
  ],
  "average_duration_seconds": 45.5
}
```

---

#### `GET /api/admin/metrics/roles`
**Descripción:** Obtiene métricas sobre roles y permisos.

**Autenticación:** Requiere rol "admin"

**Respuesta Exitosa (200):**
```json
{
  "total_roles": 5,
  "total_permissions": 20,
  "permissions_per_role": [
    {
      "role": "admin",
      "count": 15
    },
    {
      "role": "user",
      "count": 3
    }
  ]
}
```

---

#### `GET /api/admin/metrics/system`
**Descripción:** Obtiene conteos de todas las tablas importantes del sistema.

**Autenticación:** Requiere rol "admin"

**Respuesta Exitosa (200):**
```json
{
  "table_counts": {
    "users": 100,
    "active_connections": 50,
    "connection_logs": 500,
    "controls_informations": 25,
    "sessions": 200,
    "audit_runs": 200,
    "audit_script_results": 1000,
    "roles": 5,
    "permissions": 20,
    "user_roles": 150
  }
}
```

---

## 🔒 Seguridad y Autenticación

### Middleware de Autenticación
- Todos los endpoints protegidos requieren un token JWT válido en el header `Authorization: Bearer <token>`
- El middleware `RequireAuth()` valida el token y extrae información del usuario
- El middleware `RequireRole("admin")` valida que el usuario tenga el rol especificado

### Política de Sesiones
- **Sesión única:** Un usuario solo puede tener una sesión activa a la vez
- Al hacer login, si existe una sesión activa previa (y no expirada), se rechaza el nuevo login
- Las sesiones tienen una expiración de 24 horas
- El logout marca la sesión como inactiva

### Encriptación
- Las contraseñas de usuarios se hashean con bcrypt antes de almacenarse
- Las contraseñas de conexiones a bases de datos se encriptan con AES-GCM antes de almacenarse
- Las contraseñas nunca se devuelven en las respuestas de la API

---

## 📝 Notas Adicionales

1. **Gestores de Base de Datos Soportados:**
   - `pgsql` (PostgreSQL)
   - `oracle` (Oracle Database)
   - `mysql` (MySQL)
   - `mssql` (Microsoft SQL Server)
   - `otro` (Otros)

2. **Logging:**
   - Todos los endpoints registran sus acciones en logs
   - Las acciones administrativas se registran en el log de auditoría RBAC
   - Los logs incluyen información del usuario que realiza la acción

3. **Validaciones:**
   - Los endpoints validan los datos de entrada según las reglas de negocio
   - Se validan formatos de email, longitudes de campos, etc.
   - Se validan permisos y roles antes de permitir operaciones

4. **Manejo de Errores:**
   - Los errores se devuelven con códigos HTTP apropiados
   - Los mensajes de error son descriptivos pero no exponen información sensible
   - Los errores internos se registran en los logs del servidor

---

**Última actualización:** Generado automáticamente desde el análisis del código fuente.

