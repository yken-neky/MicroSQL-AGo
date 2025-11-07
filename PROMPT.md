# Prompt para GitHub Copilot - Migración a Golang con Arquitectura Hexagonal

## Contexto del Proyecto

Necesito migrar completamente un backend de Django (Python) a Golang. El sistema es una aplicación de auditoría de seguridad para servidores SQL Server que permite:

1. Gestión de usuarios con autenticación
2. Conexión a servidores SQL Server remotos
3. Ejecución de controles de auditoría predefinidos
4. Generación de reportes y logs de auditorías
5. Dashboard con métricas

**Referencia completa**: Lee detalladamente el archivo `README.md` en este mismo directorio para entender toda la funcionalidad, modelos de datos, endpoints, flujos de trabajo y detalles técnicos.

## Requisitos Principales

### Arquitectura

**IMPORTANTE**: Debes implementar **Arquitectura Hexagonal (Ports & Adapters)** con las siguientes capas:

```
backend-go/
├── cmd/
│   └── server/
│       └── main.go                    # Punto de entrada, configuración del servidor
├── internal/
│   ├── domain/                        # Capa de dominio (entidades y lógica de negocio)
│   │   ├── entities/                  # Entidades del dominio
│   │   │   ├── user.go
│   │   │   ├── connection.go
│   │   │   ├── control.go
│   │   │   └── audit.go
│   │   ├── ports/                     # Puertos (interfaces)
│   │   │   ├── repositories/          # Puertos de repositorio
│   │   │   │   ├── user_repository.go
│   │   │   │   ├── connection_repository.go
│   │   │   │   ├── control_repository.go
│   │   │   │   └── audit_repository.go
│   │   │   └── services/              # Puertos de servicios
│   │   │       ├── auth_service.go
│   │   │       ├── sql_server_service.go
│   │   │       └── audit_service.go
│   │   └── usecases/                  # Casos de uso (lógica de aplicación)
│   │       ├── user/
│   │       │   ├── register_user.go
│   │       │   ├── login_user.go
│   │       │   ├── get_profile.go
│   │       │   ├── update_profile.go
│   │       │   ├── change_password.go
│   │       │   └── deactivate_account.go
│   │       ├── connection/
│   │       │   ├── connect_to_sql_server.go
│   │       │   └── disconnect_from_sql_server.go
│   │       ├── control/
│   │       │   ├── list_controls.go
│   │       │   ├── get_control.go
│   │       │   └── execute_audit.go
│   │       └── audit/
│   │           ├── list_audits.go
│   │           ├── get_audit.go
│   │           └── get_audit_results.go
│   ├── adapters/                      # Adaptadores (implementaciones)
│   │   ├── primary/                   # Adaptadores primarios (entrada)
│   │   │   ├── http/                  # Handlers HTTP (Gin)
│   │   │   │   ├── handlers/
│   │   │   │   │   ├── user_handler.go
│   │   │   │   │   ├── connection_handler.go
│   │   │   │   │   ├── control_handler.go
│   │   │   │   │   ├── audit_handler.go
│   │   │   │   │   └── dashboard_handler.go
│   │   │   │   ├── middleware/
│   │   │   │   │   ├── auth_middleware.go
│   │   │   │   │   ├── cors_middleware.go
│   │   │   │   │   └── error_handler.go
│   │   │   │   ├── dto/               # Data Transfer Objects
│   │   │   │   │   ├── user_dto.go
│   │   │   │   │   ├── connection_dto.go
│   │   │   │   │   ├── control_dto.go
│   │   │   │   │   └── audit_dto.go
│   │   │   │   └── routes.go
│   │   │   └── cli/                   # Adaptador CLI (opcional, para futuras migraciones)
│   │   └── secondary/                 # Adaptadores secundarios (salida)
│   │       ├── persistence/           # Repositorios de base de datos
│   │       │   ├── sqlite/            # Implementación con SQLite
│   │       │   │   ├── user_repository.go
│   │       │   │   ├── connection_repository.go
│   │       │   │   ├── control_repository.go
│   │       │   │   └── audit_repository.go
│   │       │   └── migrations/        # Migraciones de base de datos
│   │       ├── external/              # Servicios externos
│   │       │   ├── sql_server/        # Cliente de SQL Server
│   │       │   │   ├── connection_pool.go
│   │       │   │   ├── query_executor.go
│   │       │   │   └── connection_manager.go
│   │       │   └── encryption/        # Servicio de encriptación
│   │       │       └── password_encryption.go
│   │       └── security/              # Servicios de seguridad
│   │           ├── jwt.go
│   │           └── password.go
│   └── config/                        # Configuración
│       ├── config.go
│       └── database.go
├── pkg/                               # Paquetes compartidos
│   └── utils/
│       ├── logger.go
│       └── errors.go
├── docs/                              # Documentación generada
│   ├── api/
│   │   └── swagger.yaml               # Documentación Swagger/OpenAPI
│   └── architecture.md                # Documentación de arquitectura
├── tests/                             # Tests
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── go.mod
├── go.sum
└── .env.example
```

### Stack Tecnológico

- **Framework Web**: Gin (github.com/gin-gonic/gin)
- **ORM**: GORM (gorm.io/gorm)
- **Base de datos**: SQLite (github.com/mattn/go-sqlite3)
- **SQL Server Driver**: github.com/denisenkom/go-mssqldb
- **Autenticación**: JWT (github.com/golang-jwt/jwt/v5)
- **Password Hashing**: golang.org/x/crypto/bcrypt
- **Validación**: github.com/go-playground/validator/v10
- **CORS**: github.com/gin-contrib/cors
- **Logging**: zap (go.uber.org/zap) o logrus
- **Documentación API**: Swagger/OpenAPI (github.com/swaggo/swag)

### Principios de Arquitectura Hexagonal

1. **Domain Layer (internal/domain/)**: 
   - Contiene las entidades puras del negocio
   - Define los puertos (interfaces) que necesita
   - NO tiene dependencias externas
   - Lógica de negocio pura

2. **Use Cases Layer (internal/domain/usecases/)**:
   - Orquesta la lógica de aplicación
   - Usa los puertos (interfaces) definidos en domain
   - Implementa los casos de uso específicos
   - NO conoce detalles de implementación

3. **Adapters Layer (internal/adapters/)**:
   - **Primary Adapters**: Implementan la entrada al sistema (HTTP handlers, CLI, etc.)
   - **Secondary Adapters**: Implementan la salida (base de datos, servicios externos, etc.)
   - Implementan las interfaces definidas en domain/ports

### Modelos de Dominio (Entidades)

Basándote en el README.md, crea las siguientes entidades:

1. **User**: Usuario del sistema
   - ID, Username, Email, Password (hashed), FirstName, LastName, Role, CreatedAt, LastLogin, IsActive

2. **ActiveConnection**: Conexión activa a SQL Server
   - ID, UserID, Driver, Server, DBUser, Password (encriptada), IsConnected, LastConnected

3. **ConnectionLog**: Log de conexiones
   - ID, UserID, Driver, Server, DBUser, Timestamp, Status

4. **ControlInformation**: Información de controles de auditoría
   - ID, Idx, Chapter, Name, Description, Impact, GoodConfig, BadConfig, Ref

5. **ControlScript**: Script SQL de un control
   - ID, ControlInfoID, ControlType (manual/automatic), QuerySQL

6. **AuditoryLog**: Registro de auditoría
   - ID, UserID, ServerID, Type (Completa/Parcial), Timestamp, Criticidad

7. **AuditoryLogResult**: Resultado individual de un control
   - ID, AuditoryLogID, ControlID, Result (TRUE/FALSE/MANUAL)

### Endpoints a Implementar

Implementa TODOS los endpoints documentados en el README.md:

#### Users API (`/api/users/`)
- POST `/api/users/register/` - Registrar usuario
- POST `/api/users/login/` - Iniciar sesión
- GET `/api/users/profile/` - Obtener perfil
- POST `/api/users/logout/` - Cerrar sesión
- PUT `/api/users/update_profile/` - Actualizar perfil
- POST `/api/users/change_password/` - Cambiar contraseña
- POST `/api/users/deactivate_account/` - Desactivar cuenta

#### Connections API (`/api/sql_conn/`)
- POST `/api/sql_conn/connections/connect/` - Conectar a SQL Server
- POST `/api/sql_conn/connections/disconnect/` - Desconectar de SQL Server
- GET `/api/sql_conn/admin/active_connections/` - Listar conexiones (admin)
- GET `/api/sql_conn/admin/active_connections/<id>/` - Detalle de conexión (admin)

#### Controls API (`/api/sql/controls/`)
- GET `/api/sql/controls/controls_info/` - Listar controles
- GET `/api/sql/controls/control_info/<id>/` - Detalle de control
- GET `/api/sql/controls/execute/` - Ejecutar auditoría (completa o parcial)

#### Logs API (`/api/logs/`)
- GET `/api/logs/connection_logs_list/` - Listar logs de conexión
- GET `/api/logs/connection_logs_details/<id>/` - Detalle de log
- GET `/api/logs/admin/connection_logs_list/` - Listar todos los logs (admin)
- GET `/api/logs/auditory_logs_list/` - Listar auditorías
- GET `/api/logs/auditory_logs_detail/<id>` - Detalle de auditoría
- GET `/api/logs/auditory_logs_results/<audit_id>/` - Resultados de auditoría
- GET `/api/logs/admin/auditory_logs_list/` - Listar todas las auditorías (admin)

#### Dashboard API (`/api/dashGET/`)
- GET `/api/dashGET/auditoryAmount/` - Cantidad de auditorías
- GET `/api/dashGET/connectionAmount/` - Cantidad de conexiones
- GET `/api/dashGET/correctRate/` - Tasa de correctitud

### Funcionalidades Críticas

1. **Conexión a SQL Server**:
   - Pool de conexiones por usuario
   - Timeout de conexión: 30 segundos
   - Reutilización de conexiones activas
   - Encriptación de contraseñas antes de almacenar
   - Validación de conexión antes de almacenar

2. **Ejecución de Auditorías**:
   - Soporte para auditoría completa (todos los controles)
   - Soporte para auditoría parcial (controles específicos por índice)
   - Manejo de controles manuales vs automáticos
   - Timeout de 30 segundos por query
   - Cálculo de criticidad automático
   - Registro de todos los resultados

3. **Autenticación y Autorización**:
   - JWT con expiración de 24 horas
   - Roles: `cliente` y `admin`
   - Middleware de autenticación
   - Validación de permisos por endpoint
   - Manejo de tokens revocados (opcional)

4. **Seguridad**:
   - Encriptación AES de contraseñas de SQL Server
   - Hash bcrypt de contraseñas de usuarios
   - Validación de entrada en todos los endpoints
   - Sanitización de queries SQL (usar parámetros preparados)
   - CORS configurado correctamente

### Requisitos de Documentación

**CRÍTICO**: Documenta TODO durante el desarrollo:

1. **Código**:
   - Comentarios en funciones públicas
   - Documentación de estructuras y tipos
   - Ejemplos de uso en comentarios
   - Documentación de errores posibles

2. **API**:
   - Swagger/OpenAPI completo con:
     - Descripción de cada endpoint
     - Parámetros de entrada
     - Respuestas posibles (éxito y error)
     - Ejemplos de requests y responses
     - Códigos de estado HTTP

3. **Arquitectura**:
   - Crea `docs/architecture.md` explicando:
     - Diagrama de capas
     - Flujo de datos
     - Decisión de diseño
     - Cómo agregar nuevos casos de uso
     - Cómo agregar nuevos adaptadores

4. **README del Proyecto**:
   - Instrucciones de instalación
   - Configuración de variables de entorno
   - Cómo ejecutar el proyecto
   - Cómo ejecutar tests
   - Cómo generar documentación Swagger

### Ejemplo de Estructura Hexagonal

#### 1. Definir Entidad (Domain)
```go
// internal/domain/entities/user.go
package entities

import "time"

type User struct {
    ID        uint
    Username  string
    Email     string
    Password  string // Hashed
    FirstName string
    LastName  string
    Role      string // "cliente" | "admin"
    CreatedAt time.Time
    LastLogin time.Time
    IsActive  bool
}
```

#### 2. Definir Puerto (Interface)
```go
// internal/domain/ports/repositories/user_repository.go
package repositories

import "your-project/internal/domain/entities"

type UserRepository interface {
    Create(user *entities.User) error
    FindByID(id uint) (*entities.User, error)
    FindByUsername(username string) (*entities.User, error)
    FindByEmail(email string) (*entities.User, error)
    Update(user *entities.User) error
    Delete(id uint) error
}
```

#### 3. Definir Caso de Uso
```go
// internal/domain/usecases/user/register_user.go
package user

import (
    "your-project/internal/domain/entities"
    "your-project/internal/domain/ports/repositories"
    "your-project/internal/domain/ports/services"
)

type RegisterUserUseCase struct {
    userRepo repositories.UserRepository
    authService services.AuthService
}

func NewRegisterUserUseCase(userRepo repositories.UserRepository, authService services.AuthService) *RegisterUserUseCase {
    return &RegisterUserUseCase{
        userRepo: userRepo,
        authService: authService,
    }
}

func (uc *RegisterUserUseCase) Execute(username, email, password string) (*entities.User, string, error) {
    // Validaciones
    // Crear usuario
    // Generar token
    // Retornar usuario y token
}
```

#### 4. Implementar Adaptador Primario (HTTP Handler)
```go
// internal/adapters/primary/http/handlers/user_handler.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "your-project/internal/domain/usecases/user"
)

type UserHandler struct {
    registerUserUseCase *user.RegisterUserUseCase
    // ... otros casos de uso
}

func (h *UserHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, token, err := h.registerUserUseCase.Execute(req.Username, req.Email, req.Password)
    // ... manejar respuesta
}
```

#### 5. Implementar Adaptador Secundario (Repository)
```go
// internal/adapters/secondary/persistence/sqlite/user_repository.go
package sqlite

import (
    "gorm.io/gorm"
    "your-project/internal/domain/entities"
    "your-project/internal/domain/ports/repositories"
)

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *entities.User) error {
    return r.db.Create(user).Error
}

// ... implementar todos los métodos de la interfaz
```

### Tests

Implementa tests para:

1. **Tests Unitarios**:
   - Casos de uso
   - Lógica de dominio
   - Servicios

2. **Tests de Integración**:
   - Repositorios con base de datos de prueba
   - Handlers HTTP con servidor de prueba
   - Conexión a SQL Server (mock o test DB)

3. **Tests E2E**:
   - Flujos completos de usuario
   - Ejecución de auditorías

### Variables de Entorno

Crea un archivo `.env.example` con:

```env
# Server
SERVER_PORT=8000
SERVER_HOST=localhost

# Database
DB_PATH=./db.sqlite3

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION_HOURS=24

# Encryption
ENCRYPTION_KEY=your-encryption-key-32-bytes

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Orden de Implementación Sugerido

1. **Fase 1: Infraestructura Base**
   - Configuración del proyecto
   - Estructura de carpetas
   - Configuración de base de datos
   - Migraciones
   - Configuración de logging

2. **Fase 2: Dominio y Casos de Uso**
   - Entidades del dominio
   - Puertos (interfaces)
   - Casos de uso de usuarios
   - Casos de uso de conexiones
   - Casos de uso de controles
   - Casos de uso de auditorías

3. **Fase 3: Adaptadores Secundarios**
   - Repositorios SQLite
   - Servicio de SQL Server
   - Servicio de encriptación
   - Servicio de JWT
   - Servicio de passwords

4. **Fase 4: Adaptadores Primarios**
   - Handlers HTTP
   - Middleware
   - DTOs
   - Rutas
   - Manejo de errores

5. **Fase 5: Funcionalidades Avanzadas**
   - Pool de conexiones
   - Timeouts
   - Validaciones
   - Métricas
   - Documentación Swagger

6. **Fase 6: Testing y Documentación**
   - Tests unitarios
   - Tests de integración
   - Tests E2E
   - Documentación completa
   - README

### Instrucciones Específicas

1. **Lee completamente el README.md** antes de empezar
2. **Sigue estrictamente la arquitectura hexagonal** - no mezcles capas
3. **Documenta cada decisión importante** en comentarios o docs
4. **Usa nombres descriptivos** en español para comentarios, en inglés para código
5. **Implementa manejo de errores robusto** con tipos de error personalizados
6. **Valida todas las entradas** de usuario
7. **Usa prepared statements** para todas las queries SQL
8. **Implementa logging estructurado** para debugging y monitoreo
9. **Crea migraciones de base de datos** para todas las tablas
10. **Genera documentación Swagger** automáticamente

### Notas Importantes

- El sistema actual almacena contraseñas de SQL Server en texto plano - **DEBES encriptarlas** en la nueva implementación
- Usa JWT stateless en lugar de tokens en base de datos (más escalable)
- Implementa pool de conexiones para SQL Server (no crear conexión por request)
- Los controles pueden ser `manual` o `automatic` - los manuales no ejecutan query SQL
- Las queries SQL retornan un solo valor que se interpreta como TRUE/FALSE
- La criticidad se calcula como porcentaje de controles fallidos
- Un usuario solo puede tener una conexión activa a la vez (OneToOne)

### Resultado Esperado

Al finalizar, deberías tener:

1. ✅ Proyecto Go completamente funcional
2. ✅ Arquitectura hexagonal bien implementada
3. ✅ Todos los endpoints funcionando
4. ✅ Tests implementados
5. ✅ Documentación Swagger completa
6. ✅ Documentación de arquitectura
7. ✅ README con instrucciones
8. ✅ Migraciones de base de datos
9. ✅ Manejo de errores robusto
10. ✅ Logging y monitoreo
11. ✅ Seguridad implementada (encriptación, JWT, validaciones)
12. ✅ Pool de conexiones a SQL Server
13. ✅ Timeouts y límites configurados

### Comenzar Implementación

Empieza por:
1. Crear la estructura de carpetas
2. Configurar `go.mod` con todas las dependencias
3. Crear archivo de configuración
4. Implementar las entidades del dominio
5. Definir los puertos (interfaces)
6. Continuar con los casos de uso y adaptadores

**Recuerda**: Documenta TODO durante el proceso. Cada función pública debe tener comentarios, cada decisión de diseño debe estar justificada, y la documentación de API debe ser completa.

