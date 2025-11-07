# SQL Server Auditing Backend - Documentación Técnica Completa

## Tabla de Contenidos
1. [Propósito del Proyecto](#propósito-del-proyecto)
2. [Arquitectura y Tecnologías](#arquitectura-y-tecnologías)
3. [Estructura del Proyecto](#estructura-del-proyecto)
4. [Modelo de Datos](#modelo-de-datos)
5. [Sistema de Conexión a SQL Server](#sistema-de-conexión-a-sql-server)
6. [Endpoints y API](#endpoints-y-api)
7. [Autenticación y Autorización](#autenticación-y-autorización)
8. [Flujos de Trabajo](#flujos-de-trabajo)
9. [Consideraciones para Migración a Golang](#consideraciones-para-migración-a-golang)
10. [Detalles de Implementación](#detalles-de-implementación)

---

## Propósito del Proyecto

Este backend es un sistema de **auditoría de seguridad y configuración para servidores SQL Server**. Permite a los usuarios:

1. **Conectarse a servidores SQL Server remotos** utilizando autenticación de SQL Server
2. **Ejecutar controles de auditoría** predefinidos que verifican configuraciones de seguridad y buenas prácticas
3. **Generar reportes de auditoría** con resultados detallados de cada control
4. **Gestionar usuarios** con roles diferenciados (Administrador y Cliente)
5. **Mantener un historial** de conexiones y auditorías realizadas

### Casos de Uso Principales

- Verificación de configuraciones de seguridad en SQL Server
- Validación de buenas prácticas de configuración
- Generación de informes de cumplimiento
- Monitoreo de cambios en la configuración del servidor
- Auditoría de múltiples servidores SQL Server desde una sola plataforma

---

## Arquitectura y Tecnologías

### Stack Tecnológico

#### Backend Framework
- **Django 5.0.10**: Framework web de Python
- **Django REST Framework 3.15.2**: Para la construcción de APIs RESTful
- **Python**: Lenguaje de programación principal

#### Base de Datos
- **SQLite3**: Base de datos local para almacenar datos de la aplicación (usuarios, conexiones, logs)
- **SQL Server** (remoto): Base de datos objetivo de las auditorías

#### Conexión a SQL Server
- **pyodbc 5.2.0**: Driver ODBC para conexión a SQL Server
- **mssql-django 1.5**: Backend de Django para SQL Server (no utilizado actualmente, solo pyodbc)
- **django-mssql-backend 2.8.1**: Dependencia adicional (no utilizada directamente)

#### Seguridad y Autenticación
- **Django REST Framework Token Authentication**: Sistema de autenticación basado en tokens
- **Django CORS Headers 4.6.0**: Manejo de CORS para comunicación con frontend
- **Django Password Validation**: Validación de contraseñas según políticas de seguridad

#### Documentación
- **CoreAPI 2.3.3**: Para documentación automática de la API

### Arquitectura de la Aplicación

El proyecto sigue la arquitectura **Django Apps** con separación de responsabilidades:

```
Backend/
├── General/              # Configuración principal del proyecto
├── Users_App/           # Gestión de usuarios y autenticación
├── Connecting_App/      # Manejo de conexiones a SQL Server
├── InsideDB_App/        # Controles de auditoría y ejecución de queries
├── Logs_App/            # Registro de logs y resultados de auditorías
└── LandingP_App/        # Landing page (vacía actualmente)
```

---

## Estructura del Proyecto

### General/
Configuración principal del proyecto Django.

**Archivos clave:**
- `settings.py`: Configuración general, apps instaladas, middleware, autenticación
- `urls.py`: Enrutamiento principal de URLs
- `wsgi.py`: Configuración WSGI para despliegue
- `asgi.py`: Configuración ASGI para despliegue asíncrono

### Users_App/
Gestión completa de usuarios, autenticación y autorización.

**Archivos:**
- `models.py`: Modelo `CustomUser` (extiende `AbstractUser`)
- `views.py`: ViewSet para registro, login, logout, perfil, cambio de contraseña
- `serializer.py`: Serializadores para validación y serialización de usuarios
- `permissions.py`: Permisos personalizados (IsAdmin, IsClient, IsOwner)
- `urls.py`: Rutas de la API de usuarios

### Connecting_App/
Manejo de conexiones activas a SQL Server.

**Archivos:**
- `models.py`: Modelo `ActiveConnection` para conexiones activas
- `views.py`: ViewSet para conectar/desconectar de SQL Server
- `serializer.py`: Serializador para ActiveConnection
- `utils.py`: Funciones de utilidad para conexión a SQL Server (pyodbc)
- `permissions.py`: Permisos relacionados con conexiones
- `urls.py`: Rutas de la API de conexiones

### InsideDB_App/
Controles de auditoría y ejecución de queries SQL.

**Archivos:**
- `models.py`: 
  - `Controls_Information`: Información de cada control de auditoría
  - `Controls_Scripts`: Scripts SQL para cada control
- `views.py`: 
  - Vistas para listar controles
  - Función `execute_query`: Ejecuta los controles de auditoría
- `serializer.py`: Serializadores para controles
- `utils.py`: Funciones para ejecutar queries y obtener conexiones activas
- `urls.py`: Rutas de la API de controles

### Logs_App/
Registro y consulta de logs de conexiones y auditorías.

**Archivos:**
- `models.py`: 
  - `ConnectionLog`: Logs de conexiones/desconexiones
  - `AuditoryLog`: Registro de auditorías ejecutadas
  - `AuditoryLogResult`: Resultados individuales de cada control en una auditoría
- `views.py`: Vistas para consultar logs y métricas del dashboard
- `serializer.py`: Serializadores para logs
- `urls.py`: Rutas de la API de logs

---

## Modelo de Datos

### Diagrama de Entidades

```
CustomUser (Users_App)
├── id (PK)
├── username
├── email
├── password (hashed)
├── first_name
├── last_name
├── role (cliente/admin)
├── created_at
├── last_login
├── is_active
└── is_superuser

ActiveConnection (Connecting_App)
├── id (PK)
├── user (FK -> CustomUser, OneToOne)
├── driver (string) - Ej: "ODBC Driver 17 for SQL Server"
├── server (string) - IP o nombre del servidor
├── db_user (string) - Usuario de SQL Server
├── password (string) - Contraseña (almacenada en texto plano)
├── is_connected (boolean)
└── last_connected (DateTime)

ConnectionLog (Logs_App)
├── id (PK)
├── user (FK -> CustomUser)
├── driver (string)
├── server (string)
├── db_user (string)
├── timestamp (DateTime)
└── status (connected/disconnected/reconnected)

Controls_Information (InsideDB_App)
├── id (PK)
├── idx (Integer) - Índice único del control
├── name (string) - Nombre del control
├── chapter (string) - Capítulo (2, 3, 4, 5, 6, 7)
├── description (Text) - Descripción del control
├── impact (Text) - Impacto de no cumplir el control
├── good_config (Text) - Mensaje cuando está bien configurado
├── bad_config (Text) - Mensaje cuando está mal configurado
└── ref (URL) - Referencia/documentación
└── UNIQUE(idx, chapter)

Controls_Scripts (InsideDB_App)
├── id (PK)
├── control_script_id (FK -> Controls_Information)
├── control_type (manual/automatic)
└── query_sql (Text) - Query SQL a ejecutar

AuditoryLog (Logs_App)
├── id (PK)
├── user (FK -> CustomUser)
├── server (FK -> ConnectionLog)
├── type (Completa/Parcial)
├── timestamp (DateTime)
└── criticidad (Float) - Porcentaje de controles fallidos

AuditoryLogResult (Logs_App)
├── id (PK)
├── auditory_log (FK -> AuditoryLog)
├── control (FK -> Controls_Information)
├── result (string) - TRUE/FALSE/MANUAL
└── UNIQUE(auditory_log, control)
```

### Modelos Detallados

#### CustomUser (Users_App/models.py)

Extiende `AbstractUser` de Django y agrega campos personalizados:

```python
class CustomUser(AbstractUser):
    created_at = models.DateTimeField(auto_now_add=True)
    last_login = models.DateTimeField(blank=True, null=True)
    role = models.CharField(max_length=10, choices=ROLE_CHOICES, default='cliente')
    
    ROLE_CHOICES = (
        ('cliente', 'Client'),
        ('admin', 'Admin'),
    )
```

**Características:**
- Sobrescribe `save()` para actualizar `last_login` automáticamente
- Roles: `cliente` (por defecto) y `admin`
- Hereda campos estándar: username, email, password, first_name, last_name, is_active, etc.

#### ActiveConnection (Connecting_App/models.py)

Almacena la conexión activa de cada usuario:

```python
class ActiveConnection(models.Model):
    user = models.OneToOneField(CustomUser, on_delete=models.CASCADE)
    driver = models.CharField(max_length=255)
    server = models.CharField(max_length=255)
    db_user = models.CharField(max_length=255)
    password = models.CharField(max_length=255)  # ⚠️ Almacenado en texto plano
    is_connected = models.BooleanField(default=False)
    last_connected = models.DateTimeField(auto_now=True)
```

**Características:**
- Relación OneToOne: un usuario solo puede tener una conexión activa a la vez
- `password` se almacena en texto plano (⚠️ **MEJORA NECESARIA**: Encriptar en producción)
- `last_connected` se actualiza automáticamente al guardar

#### Controls_Information (InsideDB_App/models.py)

Catálogo de controles de auditoría disponibles:

```python
class Controls_Information(models.Model):
    class Chapter(models.TextChoices):
        CHAPTER_TWO = '2', 'Two'
        CHAPTER_THREE = '3', 'Three'
        CHAPTER_FOUR = '4', 'Four'
        CHAPTER_FIVE = '5', 'Five'
        CHAPTER_SIX = '6', 'Six'
        CHAPTER_SEVEN = '7', 'Seven'
    
    id = models.AutoField(primary_key=True)
    idx = models.IntegerField(default=0)
    name = models.CharField(max_length=255, default="Control Name")
    chapter = models.CharField(max_length=10, choices=Chapter.choices)
    description = models.TextField(default="Description")
    impact = models.TextField(default="Impact")
    good_config = models.TextField(...)
    bad_config = models.TextField(...)
    ref = models.URLField(default="http://References.ref")
    
    class Meta:
        unique_together = [['idx', 'chapter']]
```

**Características:**
- `idx` + `chapter` forman una clave única compuesta
- Campos de texto para descripción, impacto y configuración
- URL de referencia para documentación

#### Controls_Scripts (InsideDB_App/models.py)

Scripts SQL asociados a cada control:

```python
class Controls_Scripts(models.Model):
    class ControlType(models.TextChoices):
        MANUAL = 'manual', 'Manual'
        AUTOMATIC = 'automatic', 'Automatic'
    
    control_script_id = models.ForeignKey(Controls_Information, on_delete=models.CASCADE)
    control_type = models.CharField(max_length=10, choices=ControlType.choices)
    query_sql = models.TextField(blank=True, null=True)
```

**Características:**
- `control_type`: 
  - `automatic`: Se ejecuta automáticamente con la query SQL
  - `manual`: Requiere intervención manual, no ejecuta query
- `query_sql`: Query SQL que debe retornar un valor booleano (TRUE/FALSE) o un valor que se interpreta como booleano

#### ConnectionLog (Logs_App/models.py)

Historial de conexiones y desconexiones:

```python
class ConnectionLog(models.Model):
    user = models.ForeignKey(CustomUser, on_delete=models.CASCADE)
    driver = models.CharField(max_length=255)
    server = models.CharField(max_length=255)
    db_user = models.CharField(max_length=255)
    timestamp = models.DateTimeField(auto_now_add=True)
    status = models.CharField(max_length=50, choices=[
        ('connected', 'Connected'),
        ('disconnected', 'Disconnected'),
        ('reconnected', 'Reconnected')
    ])
```

**Características:**
- Se crea un registro cada vez que un usuario se conecta o desconecta
- `timestamp` se asigna automáticamente al crear el registro

#### AuditoryLog (Logs_App/models.py)

Registro de auditorías ejecutadas:

```python
class AuditoryLog(models.Model):
    user = models.ForeignKey(CustomUser, on_delete=models.CASCADE)
    server = models.ForeignKey(ConnectionLog, on_delete=models.CASCADE)
    type = models.CharField(max_length=50, choices=[
        ('Completa', 'Completa'),
        ('parcial', 'Parcial')
    ], default='Parcial')
    timestamp = models.DateTimeField(auto_now_add=True)
    criticidad = models.FloatField(default=0)  # Porcentaje de controles fallidos
```

**Características:**
- `type`: 
  - `Completa`: Se ejecutaron todos los controles
  - `Parcial`: Se ejecutaron controles específicos (por índices)
- `criticidad`: Porcentaje calculado como `(falses / total) * 100`
- Método `calcular_criticidad()`: Calcula y actualiza el porcentaje basado en los resultados

#### AuditoryLogResult (Logs_App/models.py)

Resultados individuales de cada control en una auditoría:

```python
class AuditoryLogResult(models.Model):
    auditory_log = models.ForeignKey(AuditoryLog, on_delete=models.CASCADE, related_name='results')
    control = models.ForeignKey(Controls_Information, on_delete=models.CASCADE)
    result = models.CharField(max_length=10)  # 'TRUE', 'FALSE', 'MANUAL', etc.
    
    class Meta:
        unique_together = ('auditory_log', 'control')
```

**Características:**
- Un registro por cada control ejecutado en una auditoría
- `result`: Valores posibles:
  - `'TRUE'`: Control pasado (configuración correcta)
  - `'FALSE'`: Control fallido (configuración incorrecta)
  - `'MANUAL'`: Control manual (no se ejecuta automáticamente)
- Restricción única: un control no puede aparecer dos veces en la misma auditoría

---

## Sistema de Conexión a SQL Server

### Flujo de Conexión

1. **Usuario solicita conexión** (`POST /api/sql_conn/connections/connect/`)
2. **Validación de parámetros**: driver, server, db_user, password
3. **Verificación de conexión existente**: Se verifica si ya existe una conexión activa con los mismos parámetros
4. **Intento de conexión a SQL Server**: Se utiliza `pyodbc` para establecer la conexión
5. **Almacenamiento de conexión**: Si la conexión es exitosa, se crea/actualiza `ActiveConnection`
6. **Registro de log**: Se crea un `ConnectionLog` con status 'connected' o 'reconnected'

### Implementación Técnica

#### Función: `get_db_connection()` (Connecting_App/utils.py)

```python
def get_db_connection(driver, server, db_user, password, user):
    try:
        connection_string = f"DRIVER={{{driver}}};SERVER={server};UID={db_user};PWD={password};TrustServerCertificate=yes;"
        connection = pyodbc.connect(connection_string, autocommit=True)
        return {"status": "connected", "connection": connection}
    except pyodbc.Error as e:
        return {"status": "connection_failed", "error": str(e)}
```

**Detalles:**
- **Connection String Format**: Utiliza formato ODBC estándar
  - `DRIVER`: Nombre del driver ODBC (ej: "ODBC Driver 17 for SQL Server")
  - `SERVER`: IP o nombre del servidor SQL Server
  - `UID`: Usuario de SQL Server
  - `PWD`: Contraseña
  - `TrustServerCertificate=yes`: Confía en el certificado del servidor (útil para desarrollo, revisar en producción)
- **Autocommit**: Las transacciones se confirman automáticamente
- **Manejo de errores**: Captura `pyodbc.Error` y retorna un diccionario con el estado

#### Función: `create_or_update_connection()` (Connecting_App/utils.py)

```python
def create_or_update_connection(user, driver, server, db_user, password):
    active_conn, created = ActiveConnection.objects.update_or_create(
        user=user,
        defaults={
            'driver': driver,
            'server': server,
            'db_user': db_user,
            'password': password,
            'is_connected': True
        }
    )
    ConnectionLog.objects.create(
        user=user,
        driver=driver,
        server=server,
        db_user=db_user,
        status='connected' if created else 'reconnected'
    )
    return active_conn
```

**Detalles:**
- Utiliza `update_or_create()` para evitar duplicados
- Crea un `ConnectionLog` automáticamente
- Marca `is_connected=True`

#### Función: `obtener_conexion_activa_db()` (InsideDB_App/utils.py)

```python
def obtener_conexion_activa_db(user, connection_lock):
    with connection_lock:
        active_conn = ActiveConnection.objects.filter(user=user, is_connected=True).first()
        
        if not active_conn:
            return {"response": Response({"status": "no_connection"}, ...), "db_connection": None}
        
        connection_status = get_db_connection(
            active_conn.driver,
            active_conn.server,
            active_conn.db_user,
            active_conn.password,
            user
        )
        
        if connection_status["status"] != "connected":
            return {"response": Response({"status": "connection_failed", ...}, ...), "db_connection": None}
        
        return {"response": None, "db_connection": connection_status["connection"], "active_conn": active_conn}
```

**Detalles:**
- Utiliza un `Lock` de threading para sincronización en entornos concurrentes
- Obtiene la conexión activa del usuario desde la base de datos
- Revalida la conexión llamando a `get_db_connection()` (la conexión pyodbc puede expirar)
- Retorna el objeto `connection` de pyodbc para su uso en queries

### Gestión de Conexiones

**Características importantes:**
- ⚠️ **No hay pool de conexiones**: Cada request crea una nueva conexión
- ⚠️ **No hay cierre explícito**: Las conexiones se cierran cuando el objeto se elimina (garbage collection)
- ⚠️ **Contraseñas en texto plano**: Se almacenan sin encriptar en `ActiveConnection`
- ✅ **Sincronización con Lock**: Previene condiciones de carrera al obtener conexiones

**Mejoras recomendadas para Golang:**
- Implementar un pool de conexiones (usando `database/sql` con configuración de pool)
- Cerrar conexiones explícitamente con `defer connection.Close()`
- Encriptar contraseñas antes de almacenar (usar `crypto` de Go)
- Implementar timeouts y límites de conexión

---

## Endpoints y API

### Base URL
```
http://localhost:8000
```

### Autenticación

La mayoría de los endpoints requieren autenticación mediante **Token Authentication**.

**Headers requeridos:**
```
Authorization: Token <token_key>
```

**Obtención de token:**
- Al registrarse: `POST /api/users/register/` retorna el token
- Al hacer login: `POST /api/users/login/` retorna el token

---

### 1. Users API (`/api/users/`)

#### 1.1. Registrar Usuario
```
POST /api/users/register/
```
**Permisos:** Público (AllowAny)

**Body:**
```json
{
  "username": "usuario123",
  "email": "usuario@example.com",
  "password": "ContraseñaSegura123",
  "first_name": "Nombre",
  "last_name": "Apellido",
  "role": "cliente"  // Opcional, por defecto "cliente"
}
```

**Respuesta exitosa (201):**
```json
{
  "user": {
    "id": 1,
    "username": "usuario123",
    "email": "usuario@example.com",
    "first_name": "Nombre",
    "last_name": "Apellido",
    "created_at": "2025-01-22T10:00:00Z",
    "last_login": "2025-01-22T10:00:00Z",
    "role": "cliente"
  },
  "token": "abc123def456..."
}
```

**Cookies:**
- `authToken`: Token almacenado en cookie HttpOnly (SameSite=Lax)

**Validaciones:**
- Password debe cumplir con validadores de Django (longitud, complejidad, etc.)
- Email debe ser válido y único
- Username debe ser único

---

#### 1.2. Iniciar Sesión
```
POST /api/users/login/
```
**Permisos:** Público (AllowAny)

**Body:**
```json
{
  "username": "usuario123",
  "password": "ContraseñaSegura123"
}
```

**Respuesta exitosa (200):**
```json
{
  "token": "abc123def456...",
  "user": {
    "id": 1,
    "username": "usuario123",
    "email": "usuario@example.com",
    "first_name": "Nombre",
    "last_name": "Apellido",
    "created_at": "2025-01-22T10:00:00Z",
    "last_login": "2025-01-22T10:00:00Z",
    "role": "cliente"
  }
}
```

**Errores:**
- `400`: Credenciales inválidas
- `400`: Usuario ya está logueado (token existente)
- `403`: Usuario inactivo

---

#### 1.3. Obtener Perfil
```
GET /api/users/profile/
```
**Permisos:** IsAuthenticated

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
{
  "id": 1,
  "username": "usuario123",
  "email": "usuario@example.com",
  "first_name": "Nombre",
  "last_name": "Apellido",
  "created_at": "2025-01-22T10:00:00Z",
  "last_login": "2025-01-22T10:00:00Z",
  "role": "cliente"
}
```

---

#### 1.4. Cerrar Sesión
```
POST /api/users/logout/
```
**Permisos:** IsAuthenticated

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
{
  "message": "Logged out successfully"
}
```

**Efectos:**
- Elimina el token del usuario
- Elimina la cookie `authToken`

---

#### 1.5. Actualizar Perfil
```
PUT /api/users/update_profile/
```
**Permisos:** IsAuthenticated

**Headers:**
```
Authorization: Token <token>
```

**Body (campos opcionales):**
```json
{
  "email": "nuevo@example.com",
  "first_name": "Nuevo Nombre",
  "last_name": "Nuevo Apellido",
  "role": "admin"  // Solo si el usuario es admin
}
```

**Respuesta exitosa (200):**
```json
{
  "id": 1,
  "username": "usuario123",
  "email": "nuevo@example.com",
  "first_name": "Nuevo Nombre",
  "last_name": "Nuevo Apellido",
  "created_at": "2025-01-22T10:00:00Z",
  "last_login": "2025-01-22T10:00:00Z",
  "role": "cliente"
}
```

**Nota:** Los usuarios no-admin no pueden cambiar su `role`.

---

#### 1.6. Cambiar Contraseña
```
POST /api/users/change_password/
```
**Permisos:** IsAuthenticated

**Headers:**
```
Authorization: Token <token>
```

**Body:**
```json
{
  "current_password": "ContraseñaActual123",
  "new_password": "NuevaContraseña456",
  "confirm_password": "NuevaContraseña456"
}
```

**Respuesta exitosa (200):**
```json
{
  "message": "Password changed successfully"
}
```

**Validaciones:**
- `current_password` debe ser correcta
- `new_password` debe cumplir con validadores de Django
- `new_password` y `confirm_password` deben coincidir

---

#### 1.7. Desactivar Cuenta
```
POST /api/users/deactivate_account/
```
**Permisos:** IsAuthenticated

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
{
  "message": "Account deactivated and logged out successfully"
}
```

**Efectos:**
- Establece `is_active=False` en el usuario
- Elimina el token
- Elimina la cookie `authToken`
- El usuario no podrá iniciar sesión hasta que se reactive la cuenta

---

### 2. Connections API (`/api/sql_conn/`)

#### 2.1. Conectar a SQL Server
```
POST /api/sql_conn/connections/connect/
```
**Permisos:** IsAuthenticated, IsClient

**Headers:**
```
Authorization: Token <token>
```

**Body:**
```json
{
  "driver": "ODBC Driver 17 for SQL Server",
  "server": "192.168.1.100",
  "db_user": "sa",
  "password": "password123"
}
```

**Parámetros:**
- `driver`: Nombre del driver ODBC instalado en el servidor
  - Ejemplos comunes:
    - `ODBC Driver 17 for SQL Server`
    - `ODBC Driver 13 for SQL Server`
    - `SQL Server Native Client 11.0`
- `server`: IP o nombre del servidor SQL Server (ej: `192.168.1.100` o `SQLSERVER01`)
- `db_user`: Usuario de SQL Server (ej: `sa`, `admin`, etc.)
- `password`: Contraseña del usuario de SQL Server

**Respuesta exitosa (200):**
```json
{
  "message": "Conexión creada o actualizada.",
  "connection": 1
}
```

**Errores:**
- `400`: Campos faltantes
- `400`: Ya existe una conexión activa con esos parámetros
- `500`: Error de conexión a SQL Server (credenciales incorrectas, servidor inaccesible, etc.)

**Flujo:**
1. Valida que todos los campos estén presentes
2. Verifica si ya existe una conexión activa con los mismos parámetros
3. Intenta conectar a SQL Server usando `pyodbc`
4. Si la conexión es exitosa, crea/actualiza `ActiveConnection`
5. Crea un `ConnectionLog` con status 'connected' o 'reconnected'

---

#### 2.2. Desconectar de SQL Server
```
POST /api/sql_conn/connections/disconnect/
```
**Permisos:** IsAuthenticated, IsClient

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
{
  "message": "Desconexión exitosa."
}
```

**Errores:**
- `400`: No hay conexiones activas para desconectar

**Flujo:**
1. Busca la conexión activa del usuario
2. Establece `is_connected=False`
3. Crea un `ConnectionLog` con status 'disconnected'

---

#### 2.3. Listar Todas las Conexiones Activas (Admin)
```
GET /api/sql_conn/admin/active_connections/
```
**Permisos:** IsAdmin

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
[
  {
    "id": 1,
    "user": 1,
    "driver": "ODBC Driver 17 for SQL Server",
    "server": "192.168.1.100",
    "db_user": "sa",
    "password": "password123",
    "is_connected": true,
    "last_connected": "2025-01-22T10:00:00Z"
  }
]
```

---

#### 2.4. Obtener Detalle de Conexión Activa (Admin)
```
GET /api/sql_conn/admin/active_connections/<id>/
```
**Permisos:** IsAdmin

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
{
  "id": 1,
  "user": 1,
  "driver": "ODBC Driver 17 for SQL Server",
  "server": "192.168.1.100",
  "db_user": "sa",
  "password": "password123",
  "is_connected": true,
  "last_connected": "2025-01-22T10:00:00Z"
}
```

---

### 3. Controls API (`/api/sql/controls/`)

#### 3.1. Listar Todos los Controles
```
GET /api/sql/controls/controls_info/
```
**Permisos:** IsAuthenticated

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
[
  {
    "id": 1,
    "idx": 1,
    "name": "Control de autenticación mixta",
    "chapter": "2",
    "description": "Verifica si la autenticación mixta está habilitada",
    "impact": "Riesgo de seguridad si está deshabilitada",
    "good_config": "Este parámetro se encuentra configurado de forma correcta...",
    "bad_config": "Este parámetro se encuentra configurado de forma incorrecta...",
    "ref": "https://docs.microsoft.com/..."
  }
]
```

---

#### 3.2. Obtener Detalle de un Control
```
GET /api/sql/controls/control_info/<id>/
```
**Permisos:** IsAuthenticated

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
{
  "id": 1,
  "idx": 1,
  "name": "Control de autenticación mixta",
  "chapter": "2",
  "description": "Verifica si la autenticación mixta está habilitada",
  "impact": "Riesgo de seguridad si está deshabilitada",
  "good_config": "Este parámetro se encuentra configurado de forma correcta...",
  "bad_config": "Este parámetro se encuentra configurado de forma incorrecta...",
  "ref": "https://docs.microsoft.com/..."
}
```

---

#### 3.3. Ejecutar Controles de Auditoría
```
GET /api/sql/controls/execute/
```
**Permisos:** IsAuthenticated, IsClient

**Headers:**
```
Authorization: Token <token>
```

**Query Parameters:**
- Sin parámetros: Ejecuta **todos** los controles (auditoría completa)
- Con parámetros: Ejecuta controles específicos (auditoría parcial)
  ```
  ?idxes=1&idxes=3&idxes=5
  ```

**Ejemplos:**
- Auditoría completa: `GET /api/sql/controls/execute/`
- Auditoría parcial: `GET /api/sql/controls/execute/?idxes=1&idxes=3&idxes=5`

**Respuesta exitosa (200):**
```json
{
  "status": "queries_executed",
  "control_results": {
    "1": "TRUE",
    "2": "FALSE",
    "3": "MANUAL",
    "4": "TRUE"
  },
  "audit_id": 42
}
```

**Flujo detallado:**
1. Obtiene la conexión activa del usuario usando `obtener_conexion_activa_db()`
2. Si no hay conexión activa, retorna error 400
3. Obtiene los controles a ejecutar:
   - Si hay `idxes` en query params: filtra por esos índices
   - Si no hay `idxes`: obtiene todos los controles
4. Para cada control:
   - Si es `manual`: asigna resultado `"MANUAL"`
   - Si es `automatic`: ejecuta `query_sql` usando `ejecutar_consulta()`
5. Crea un `AuditoryLog` con type 'Completa' o 'Parcial'
6. Para cada resultado, crea un `AuditoryLogResult`
7. Calcula `criticidad` como porcentaje de controles fallidos (`FALSE`)
8. Retorna los resultados y el `audit_id`

**Errores:**
- `400`: No hay conexión activa
- `404`: No se encontraron controles con los índices proporcionados
- `500`: Error al ejecutar query SQL
- `500`: Error de conexión a SQL Server

**Nota importante:** Las queries SQL deben retornar un valor que se pueda interpretar como booleano. El sistema espera que el primer valor de la primera fila sea `TRUE` o `FALSE` (o un valor que se pueda convertir a booleano).

---

#### 3.4. Listar Scripts de Controles (Admin)
```
GET /api/sql/controls/control_scripts/  # (Si existe, verificar en código)
```
**Permisos:** IsAdmin

**Nota:** Este endpoint puede no estar expuesto en las URLs. Los scripts generalmente se gestionan desde el admin de Django.

---

### 4. Logs API (`/api/logs/`)

#### 4.1. Listar Logs de Conexión del Usuario
```
GET /api/logs/connection_logs_list/
```
**Permisos:** IsClientAndOwner (solo los propios logs)

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
[
  {
    "id": 1,
    "user": 1,
    "driver": "ODBC Driver 17 for SQL Server",
    "server": "192.168.1.100",
    "db_user": "sa",
    "timestamp": "2025-01-22T10:00:00Z",
    "status": "connected"
  },
  {
    "id": 2,
    "user": 1,
    "driver": "ODBC Driver 17 for SQL Server",
    "server": "192.168.1.100",
    "db_user": "sa",
    "timestamp": "2025-01-22T11:00:00Z",
    "status": "disconnected"
  }
]
```

---

#### 4.2. Obtener Detalle de Log de Conexión
```
GET /api/logs/connection_logs_details/<id>/
```
**Permisos:** IsClient (solo los propios logs)

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
{
  "id": 1,
  "user": 1,
  "driver": "ODBC Driver 17 for SQL Server",
  "server": "192.168.1.100",
  "db_user": "sa",
  "timestamp": "2025-01-22T10:00:00Z",
  "status": "connected"
}
```

---

#### 4.3. Listar Todas las Auditorías (Admin)
```
GET /api/logs/admin/auditory_logs_list/
```
**Permisos:** IsAdmin

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
[
  {
    "id": 1,
    "user": "usuario123",
    "server": "192.168.1.100",
    "type": "Completa",
    "timestamp": "2025-01-22T10:00:00Z",
    "criticidad": 25.5,
    "results": [
      {
        "id": 1,
        "control": 1,
        "result": "TRUE"
      },
      {
        "id": 2,
        "control": 2,
        "result": "FALSE"
      }
    ]
  }
]
```

---

#### 4.4. Listar Auditorías del Usuario
```
GET /api/logs/auditory_logs_list/
```
**Permisos:** IsClient (solo las propias auditorías)

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
[
  {
    "id": 1,
    "user": "usuario123",
    "server": "192.168.1.100",
    "type": "Completa",
    "timestamp": "2025-01-22T10:00:00Z",
    "criticidad": 25.5,
    "results": [
      {
        "id": 1,
        "control": 1,
        "result": "TRUE"
      },
      {
        "id": 2,
        "control": 2,
        "result": "FALSE"
      }
    ]
  }
]
```

---

#### 4.5. Obtener Detalle de Auditoría
```
GET /api/logs/auditory_logs_detail/<id>
```
**Permisos:** IsClientAndOwner (solo las propias auditorías)

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
{
  "id": 1,
  "user": "usuario123",
  "server": "192.168.1.100",
  "type": "Completa",
  "timestamp": "2025-01-22T10:00:00Z",
  "criticidad": 25.5,
  "results": [
    {
      "id": 1,
      "control": 1,
      "result": "TRUE"
    },
    {
      "id": 2,
      "control": 2,
      "result": "FALSE"
    }
  ]
}
```

---

#### 4.6. Obtener Resultados de una Auditoría
```
GET /api/logs/auditory_logs_results/<audit_id>/
```
**Permisos:** IsClient (solo las propias auditorías)

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
[
  {
    "id": 1,
    "control": 1,
    "result": "TRUE"
  },
  {
    "id": 2,
    "control": 2,
    "result": "FALSE"
  }
]
```

**Errores:**
- `404`: Auditoría no encontrada o no pertenece al usuario

---

### 5. Dashboard API (`/api/dashGET/`)

#### 5.1. Cantidad de Auditorías
```
GET /api/dashGET/auditoryAmount/
```
**Permisos:** IsAuthenticated

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
{
  "auditTotal": 15
}
```

---

#### 5.2. Cantidad de Conexiones
```
GET /api/dashGET/connectionAmount/
```
**Permisos:** IsAuthenticated

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
{
  "connectionTotal": 8
}
```

---

#### 5.3. Tasa de Correctitud (Criticidad)
```
GET /api/dashGET/correctRate/
```
**Permisos:** IsAuthenticated

**Headers:**
```
Authorization: Token <token>
```

**Respuesta exitosa (200):**
```json
{
  "percentage": 25.5
}
```

**Nota:** Retorna la `criticidad` de la última auditoría **completa** del usuario. Si no hay auditorías completas, retorna `0`.

---

## Autenticación y Autorización

### Sistema de Autenticación

El sistema utiliza **Django REST Framework Token Authentication**.

#### Obtención de Token

1. **Registro**: Al registrar un usuario, se genera un token automáticamente
2. **Login**: Al hacer login, se genera/recupera un token existente
3. **Logout**: Al hacer logout, se elimina el token

#### Uso del Token

Los tokens se envían en el header `Authorization`:
```
Authorization: Token <token_key>
```

#### Almacenamiento de Tokens

- **Base de datos**: Tabla `authtoken_token` (generada por Django REST Framework)
- **Cookies**: Cookie `authToken` (HttpOnly, SameSite=Lax) se establece en registro/login

### Sistema de Permisos

#### Permisos Estándar de Django REST Framework

- **AllowAny**: Acceso público (sin autenticación)
- **IsAuthenticated**: Requiere autenticación (token válido)

#### Permisos Personalizados

##### IsAdmin (Users_App/permissions.py)
```python
class IsAdmin(BasePermission):
    def has_permission(self, request, view):
        return request.user.role == 'admin'
```
**Uso:** Solo usuarios con `role='admin'` pueden acceder.

##### IsClient (Users_App/permissions.py)
```python
class IsClient(BasePermission):
    def has_permission(self, request, view):
        return request.user.role == 'cliente'
```
**Uso:** Solo usuarios con `role='cliente'` pueden acceder.

##### IsClientAndOwner (Users_App/permissions.py)
```python
class IsClientAndOwner(BasePermission):
    def has_object_permission(self, request, view, obj):
        return request.user.role == 'cliente' and obj.user == request.user
```
**Uso:** Solo usuarios cliente que sean dueños del objeto pueden acceder.

##### IsOwner (Users_App/permissions.py)
```python
class IsOwner(BasePermission):
    def has_object_permission(self, request, view, obj):
        return obj.user == request.user
```
**Uso:** Solo el dueño del objeto puede acceder (independientemente del rol).

##### HasOnServiceCookie (Connecting_App/permissions.py)
```python
class HasOnServiceCookie(BasePermission):
    def has_permission(self, request, view):
        user = request.user
        active_conn = ActiveConnection.objects.filter(user=user).first()
        if active_conn.is_connected == True:
            return active_conn
```
**Uso:** Verifica que el usuario tenga una conexión activa. (⚠️ Nota: Este permiso no está siendo usado actualmente en el código)

##### NoOnServiceAccess (Connecting_App/permissions.py)
```python
class NoOnServiceAccess(BasePermission):
    def has_permission(self, request, view):
        user = request.user
        active_conn = ActiveConnection.objects.filter(user=user).first()
        if active_conn.is_connected != True:
            return active_conn
```
**Uso:** Verifica que el usuario NO tenga una conexión activa. (⚠️ Usado en UserViewSet pero la lógica parece invertida)

### Matriz de Permisos por Endpoint

| Endpoint | Permisos Requeridos |
|----------|---------------------|
| `POST /api/users/register/` | AllowAny |
| `POST /api/users/login/` | AllowAny |
| `GET /api/users/profile/` | IsAuthenticated |
| `POST /api/users/logout/` | IsAuthenticated |
| `PUT /api/users/update_profile/` | IsAuthenticated |
| `POST /api/users/change_password/` | IsAuthenticated |
| `POST /api/users/deactivate_account/` | IsAuthenticated |
| `POST /api/sql_conn/connections/connect/` | IsAuthenticated, IsClient |
| `POST /api/sql_conn/connections/disconnect/` | IsAuthenticated, IsClient |
| `GET /api/sql_conn/admin/active_connections/` | IsAdmin |
| `GET /api/sql/controls/controls_info/` | IsAuthenticated |
| `GET /api/sql/controls/control_info/<id>/` | IsAuthenticated |
| `GET /api/sql/controls/execute/` | IsAuthenticated, IsClient |
| `GET /api/logs/connection_logs_list/` | IsClientAndOwner |
| `GET /api/logs/auditory_logs_list/` | IsClient |
| `GET /api/logs/auditory_logs_detail/<id>` | IsClientAndOwner |
| `GET /api/dashGET/auditoryAmount/` | IsAuthenticated |
| `GET /api/dashGET/connectionAmount/` | IsAuthenticated |
| `GET /api/dashGET/correctRate/` | IsAuthenticated |

---

## Flujos de Trabajo

### Flujo 1: Registro y Primera Conexión

1. **Usuario se registra**
   - `POST /api/users/register/`
   - Recibe token de autenticación
   - Usuario creado con `role='cliente'` por defecto

2. **Usuario se conecta a SQL Server**
   - `POST /api/sql_conn/connections/connect/`
   - Envía: driver, server, db_user, password
   - Sistema valida conexión con `pyodbc`
   - Se crea `ActiveConnection` con `is_connected=True`
   - Se crea `ConnectionLog` con status 'connected'

3. **Usuario ejecuta auditoría**
   - `GET /api/sql/controls/execute/`
   - Sistema obtiene conexión activa
   - Ejecuta queries SQL de controles automáticos
   - Crea `AuditoryLog` y `AuditoryLogResult` para cada control
   - Calcula `criticidad`

4. **Usuario consulta resultados**
   - `GET /api/logs/auditory_logs_list/`
   - `GET /api/logs/auditory_logs_detail/<id>`
   - `GET /api/logs/auditory_logs_results/<audit_id>/`

5. **Usuario se desconecta**
   - `POST /api/sql_conn/connections/disconnect/`
   - Se actualiza `ActiveConnection` con `is_connected=False`
   - Se crea `ConnectionLog` con status 'disconnected'

### Flujo 2: Ejecución de Auditoría Parcial

1. **Usuario consulta controles disponibles**
   - `GET /api/sql/controls/controls_info/`
   - Ve lista de controles con índices y capítulos

2. **Usuario selecciona controles específicos**
   - `GET /api/sql/controls/execute/?idxes=1&idxes=3&idxes=5`
   - Sistema filtra controles por índices
   - Ejecuta solo esos controles
   - Crea `AuditoryLog` con type='Parcial'

3. **Usuario consulta resultados**
   - Similar al flujo anterior

### Flujo 3: Dashboard del Usuario

1. **Usuario consulta métricas**
   - `GET /api/dashGET/auditoryAmount/` - Total de auditorías
   - `GET /api/dashGET/connectionAmount/` - Total de conexiones
   - `GET /api/dashGET/correctRate/` - Criticidad de última auditoría completa

2. **Usuario consulta historial**
   - `GET /api/logs/connection_logs_list/` - Historial de conexiones
   - `GET /api/logs/auditory_logs_list/` - Historial de auditorías

---

## Consideraciones para Migración a Golang

### 1. Framework y Librerías

#### Framework Web
- **Recomendado**: `Gin` o `Echo` (frameworks ligeros y rápidos)
- **Alternativa**: `Fiber` (inspirado en Express.js)
- **Estándar**: `net/http` (más verboso pero sin dependencias)

#### Base de Datos
- **SQLite**: `github.com/mattn/go-sqlite3` (driver para SQLite)
- **ORM**: `gorm.io/gorm` (ORM popular) o `ent` (Facebook)
- **Migrations**: `golang-migrate/migrate` o `gorm` migrations

#### Conexión a SQL Server
- **Driver**: `github.com/denisenkom/go-mssqldb` (driver oficial para SQL Server)
- **Pool de conexiones**: Usar `database/sql` con configuración de pool
- **Connection String**: Similar a pyodbc pero formato Go:
  ```go
  connString := fmt.Sprintf("server=%s;user id=%s;password=%s;trustservercertificate=true", server, user, password)
  ```

#### Autenticación
- **JWT**: `github.com/golang-jwt/jwt/v5` (más moderno que tokens de Django)
- **Sessions**: `github.com/gorilla/sessions` (si se mantiene sesiones)
- **Password Hashing**: `golang.org/x/crypto/bcrypt` (similar a Django)

#### Validación
- **Struct Validation**: `github.com/go-playground/validator/v10`
- **Request Parsing**: `github.com/gin-gonic/gin` (binding automático)

#### CORS
- **Middleware**: `github.com/gin-contrib/cors` (para Gin)
- O implementar manualmente con headers

### 2. Estructura del Proyecto

```
backend-go/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   ├── user.go
│   │   ├── connection.go
│   │   ├── control.go
│   │   └── log.go
│   ├── handlers/
│   │   ├── user.go
│   │   ├── connection.go
│   │   ├── control.go
│   │   └── log.go
│   ├── services/
│   │   ├── auth.go
│   │   ├── sql_server.go
│   │   └── audit.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── cors.go
│   └── database/
│       ├── sqlite.go
│       └── migrations/
├── pkg/
│   └── utils/
│       └── password.go
└── go.mod
```

### 3. Modelos de Datos en Go

#### CustomUser
```go
type CustomUser struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"unique;not null" json:"username"`
    Email     string    `gorm:"unique;not null" json:"email"`
    Password  string    `gorm:"not null" json:"-"` // No se serializa
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Role      string    `gorm:"default:cliente" json:"role"` // cliente, admin
    CreatedAt time.Time `json:"created_at"`
    LastLogin time.Time `json:"last_login"`
    IsActive  bool      `gorm:"default:true" json:"is_active"`
}
```

#### ActiveConnection
```go
type ActiveConnection struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    UserID       uint      `gorm:"uniqueIndex;not null" json:"user_id"`
    User         CustomUser `gorm:"foreignKey:UserID" json:"-"`
    Driver       string    `json:"driver"`
    Server       string    `json:"server"`
    DBUser       string    `gorm:"column:db_user" json:"db_user"`
    Password     string    `json:"-"` // ⚠️ Encriptar en producción
    IsConnected  bool      `gorm:"default:false" json:"is_connected"`
    LastConnected time.Time `json:"last_connected"`
}
```

#### Controls_Information
```go
type ControlsInformation struct {
    ID         uint   `gorm:"primaryKey" json:"id"`
    Idx        int    `gorm:"uniqueIndex:idx_chapter" json:"idx"`
    Chapter    string `gorm:"uniqueIndex:idx_chapter" json:"chapter"` // 2, 3, 4, 5, 6, 7
    Name       string `json:"name"`
    Description string `gorm:"type:text" json:"description"`
    Impact     string `gorm:"type:text" json:"impact"`
    GoodConfig string `gorm:"type:text" json:"good_config"`
    BadConfig  string `gorm:"type:text" json:"bad_config"`
    Ref        string `gorm:"type:varchar(500)" json:"ref"`
}
```

#### Controls_Scripts
```go
type ControlsScripts struct {
    ID             uint              `gorm:"primaryKey" json:"id"`
    ControlInfoID  uint              `gorm:"not null" json:"control_script_id"`
    ControlInfo    ControlsInformation `gorm:"foreignKey:ControlInfoID" json:"-"`
    ControlType    string            `gorm:"default:automatic" json:"control_type"` // manual, automatic
    QuerySQL       string            `gorm:"type:text" json:"query_sql"`
}
```

#### ConnectionLog
```go
type ConnectionLog struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null" json:"user_id"`
    User      CustomUser `gorm:"foreignKey:UserID" json:"-"`
    Driver    string    `json:"driver"`
    Server    string    `json:"server"`
    DBUser    string    `gorm:"column:db_user" json:"db_user"`
    Timestamp time.Time `json:"timestamp"`
    Status    string    `json:"status"` // connected, disconnected, reconnected
}
```

#### AuditoryLog
```go
type AuditoryLog struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    UserID     uint      `gorm:"not null" json:"user_id"`
    User       CustomUser `gorm:"foreignKey:UserID" json:"-"`
    ServerID   uint      `gorm:"not null" json:"server_id"`
    Server     ConnectionLog `gorm:"foreignKey:ServerID" json:"-"`
    Type       string    `gorm:"default:parcial" json:"type"` // Completa, Parcial
    Timestamp  time.Time `json:"timestamp"`
    Criticidad float64   `gorm:"default:0" json:"criticidad"`
    Results    []AuditoryLogResult `gorm:"foreignKey:AuditoryLogID" json:"results"`
}
```

#### AuditoryLogResult
```go
type AuditoryLogResult struct {
    ID            uint              `gorm:"primaryKey" json:"id"`
    AuditoryLogID uint              `gorm:"uniqueIndex:audit_control" json:"auditory_log_id"`
    AuditoryLog   AuditoryLog       `gorm:"foreignKey:AuditoryLogID" json:"-"`
    ControlID     uint              `gorm:"uniqueIndex:audit_control" json:"control_id"`
    Control       ControlsInformation `gorm:"foreignKey:ControlID" json:"-"`
    Result        string            `json:"result"` // TRUE, FALSE, MANUAL
}
```

### 4. Conexión a SQL Server en Go

#### Función de Conexión
```go
package services

import (
    "context"
    "database/sql"
    "fmt"
    "sync"
    "time"
    _ "github.com/denisenkom/go-mssqldb"
    "your-project/models"
)

var (
    connectionPools map[uint]*sql.DB
    poolMutex       sync.RWMutex
)

func init() {
    connectionPools = make(map[uint]*sql.DB)
}

// GetDBConnection establece una conexión a SQL Server
func GetDBConnection(driver, server, dbUser, password string) (*sql.DB, error) {
    // Nota: El driver de Go no usa el parámetro "driver" como pyodbc
    // Se conecta directamente usando el driver go-mssqldb
    connString := fmt.Sprintf(
        "server=%s;user id=%s;password=%s;trustservercertificate=true;connection timeout=30",
        server, dbUser, password,
    )
    
    db, err := sql.Open("sqlserver", connString)
    if err != nil {
        return nil, fmt.Errorf("error opening connection: %w", err)
    }
    
    // Configurar pool de conexiones
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Verificar la conexión
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        db.Close()
        return nil, fmt.Errorf("error pinging database: %w", err)
    }
    
    return db, nil
}

// ObtenerConexionActivaDB obtiene la conexión activa de un usuario
func ObtenerConexionActivaDB(userID uint, activeConn *models.ActiveConnection) (*sql.DB, error) {
    poolMutex.RLock()
    if db, exists := connectionPools[userID]; exists {
        // Verificar que la conexión aún es válida
        if err := db.Ping(); err == nil {
            poolMutex.RUnlock()
            return db, nil
        }
        // Conexión expirada, eliminarla del pool
        db.Close()
        delete(connectionPools, userID)
    }
    poolMutex.RUnlock()
    
    // Crear nueva conexión
    db, err := GetDBConnection(
        activeConn.Driver,
        activeConn.Server,
        activeConn.DBUser,
        activeConn.Password,
    )
    if err != nil {
        return nil, err
    }
    
    // Almacenar en el pool
    poolMutex.Lock()
    connectionPools[userID] = db
    poolMutex.Unlock()
    
    return db, nil
}
```

#### Ejecutar Query
```go
// EjecutarConsulta ejecuta una query SQL y retorna el resultado
func EjecutarConsulta(db *sql.DB, querySQL string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    var result interface{}
    err := db.QueryRowContext(ctx, querySQL).Scan(&result)
    if err != nil {
        if err == sql.ErrNoRows {
            return "FALSE", nil
        }
        return "", fmt.Errorf("error executing query: %w", err)
    }
    
    // Convertir resultado a string TRUE/FALSE
    resultStr := fmt.Sprintf("%v", result)
    
    // Interpretar el resultado (puede ser booleano, bit, int, etc.)
    switch v := result.(type) {
    case bool:
        if v {
            return "TRUE", nil
        }
        return "FALSE", nil
    case int:
        if v != 0 {
            return "TRUE", nil
        }
        return "FALSE", nil
    case int64:
        if v != 0 {
            return "TRUE", nil
        }
        return "FALSE", nil
    case string:
        if strings.ToUpper(v) == "TRUE" || v == "1" {
            return "TRUE", nil
        }
        return "FALSE", nil
    default:
        // Intentar convertir a string y verificar
        if strings.Contains(strings.ToUpper(resultStr), "TRUE") || resultStr == "1" {
            return "TRUE", nil
        }
        return "FALSE", nil
    }
}
```

**Nota**: Se necesita importar `strings`:
```go
import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "sync"
    "time"
    _ "github.com/denisenkom/go-mssqldb"
    "your-project/models"
)
```

### 5. Autenticación en Go

#### JWT Token Authentication
```go
package services

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "your-project/models"
)

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key") // ⚠️ Usar variable de entorno en producción

// GenerateToken genera un JWT token para un usuario
func GenerateToken(user *models.CustomUser) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    
    claims := &Claims{
        UserID:   user.ID,
        Username: user.Username,
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }
    
    return tokenString, nil
}

// ValidateToken valida un JWT token
func ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    
    return claims, nil
}
```

#### Password Hashing
```go
package utils

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword hashea una contraseña usando bcrypt
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPassword verifica una contraseña contra un hash
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 6. Handlers en Go (Ejemplo con Gin)

#### Handler de Conexión
```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type ConnectionHandler struct {
    connectionService *services.ConnectionService
}

func (h *ConnectionHandler) Connect(c *gin.Context) {
    var req struct {
        Driver   string `json:"driver" binding:"required"`
        Server   string `json:"server" binding:"required"`
        DBUser   string `json:"db_user" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    userID := c.GetUint("user_id") // Obtenido del middleware de autenticación
    
    result, err := h.connectionService.Connect(userID, req.Driver, req.Server, req.DBUser, req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message":    "Conexión creada o actualizada.",
        "connection": result.ID,
    })
}

func (h *ConnectionHandler) Disconnect(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    if err := h.connectionService.Disconnect(userID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Desconexión exitosa."})
}
```

#### Handler de Ejecución de Auditoría
```go
func (h *ControlHandler) ExecuteQuery(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    // Obtener índices de query params
    idxes := c.QueryArray("idxes")
    
    results, auditID, err := h.controlService.ExecuteAudit(userID, idxes)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status":         "queries_executed",
        "control_results": results,
        "audit_id":       auditID,
    })
}
```

### 7. Middleware de Autenticación
```go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "your-project/services"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Formato: "Token <token>" o "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        claims, err := services.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // Agregar información del usuario al contexto
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        
        c.Next()
    }
}
```

### 8. Migraciones de Base de Datos

#### Usando GORM Migrations
```go
package database

import (
    "gorm.io/gorm"
    "your-project/models"
)

func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.CustomUser{},
        &models.ActiveConnection{},
        &models.ControlsInformation{},
        &models.ControlsScripts{},
        &models.ConnectionLog{},
        &models.AuditoryLog{},
        &models.AuditoryLogResult{},
    )
}
```

### 9. Configuración de CORS
```go
package main

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "time"
)

func setupCORS(router *gin.Engine) {
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
}
```

### 10. Puntos Clave de la Migración

#### Diferencias Importantes

1. **Gestión de Conexiones**:
   - En Django: Las conexiones se crean por request y se cierran automáticamente
   - En Go: Implementar pool de conexiones explícito y cerrar conexiones con `defer`

2. **Autenticación**:
   - En Django: Tokens almacenados en base de datos
   - En Go: JWT stateless (más escalable) o tokens en base de datos si se requiere revocación

3. **Validación**:
   - En Django: Serializers automáticos
   - En Go: Usar `validator` package o validación manual

4. **ORM**:
   - En Django: Django ORM (muy automático)
   - En Go: GORM requiere más configuración manual pero más control

5. **Manejo de Errores**:
   - En Django: Excepciones
   - En Go: Retornar errores explícitamente

6. **Concurrencia**:
   - En Django: GIL de Python limita verdadera concurrencia
   - En Go: Goroutines nativas para mejor rendimiento

---

## Detalles de Implementación

### 1. Ejecución de Queries SQL

#### Formato de Queries Esperadas

Las queries SQL almacenadas en `Controls_Scripts.query_sql` deben cumplir con las siguientes características:

1. **Retorno de un solo valor**: La query debe retornar un solo valor (primera fila, primera columna)
2. **Interpretación booleana**: El valor se interpreta como:
   - `TRUE`: Control pasado (configuración correcta)
   - `FALSE`: Control fallido (configuración incorrecta)
   - `MANUAL`: No se ejecuta automáticamente

#### Ejemplo de Query SQL
```sql
-- Ejemplo: Verificar si la autenticación mixta está habilitada
SELECT CASE 
    WHEN value = 1 THEN 'TRUE'
    ELSE 'FALSE'
END as result
FROM sys.configurations
WHERE name = 'user options'
```

#### Manejo de Resultados

El sistema actual en Python:
```python
def ejecutar_consulta(cursor, consulta_sql):
    cursor.execute(consulta_sql)
    filas = cursor.fetchall()
    if filas:
        resultado = filas[0][0]  # Primera fila, primera columna
        return {"status": "success", "data": resultado}
    else:
        return {"status": "success", "data": 'FALSE'}
```

**En Golang**, se debe implementar similar:
- Si la query retorna filas: tomar el primer valor
- Si no retorna filas: retornar `"FALSE"`
- Convertir el valor a string `"TRUE"` o `"FALSE"` según la lógica del control

### 2. Cálculo de Criticidad

La criticidad se calcula como el porcentaje de controles fallidos:

```python
falses = 0
total = 0
for idx, result in control_results.items():
    if result == 'FALSE':
        falses += 1
    total += 1
criticidad = round((falses / total) * 100, 2) if total > 0 else 0
```

**En Golang**:
```go
import (
    "math"
)

func CalculateCriticidad(results map[int]string) float64 {
    total := len(results)
    if total == 0 {
        return 0
    }
    
    falses := 0
    for _, result := range results {
        if result == "FALSE" {
            falses++
        }
    }
    
    return math.Round((float64(falses)/float64(total))*100*100) / 100
}
```

### 3. Gestión de Contraseñas

⚠️ **PROBLEMA DE SEGURIDAD ACTUAL**: Las contraseñas de SQL Server se almacenan en texto plano en `ActiveConnection.password`.

**Mejoras recomendadas para Golang**:

1. **Encriptar contraseñas antes de almacenar**:
```go
package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

func EncryptPassword(password string, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    ciphertext := make([]byte, aes.BlockSize+len(password))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }
    
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(password))
    
    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptPassword(encryptedPassword string, key []byte) (string, error) {
    ciphertext, err := base64.URLEncoding.DecodeString(encryptedPassword)
    if err != nil {
        return "", err
    }
    
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    if len(ciphertext) < aes.BlockSize {
        return "", fmt.Errorf("ciphertext too short")
    }
    
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]
    
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)
    
    return string(ciphertext), nil
}
```

2. **Usar variables de entorno para la clave de encriptación**
3. **Considerar usar un servicio de gestión de secretos** (HashiCorp Vault, AWS Secrets Manager, etc.)

### 4. Sincronización de Conexiones

El sistema actual usa un `Lock` de threading para prevenir condiciones de carrera:

```python
connection_lock = Lock()

def obtener_conexion_activa_db(user, connection_lock):
    with connection_lock:
        # ... código de obtención de conexión
```

**En Golang**, usar `sync.RWMutex` o `sync.Mutex`:
```go
var (
    connectionMutex sync.RWMutex
    connections     map[uint]*sql.DB
)

func GetActiveConnection(userID uint) (*sql.DB, error) {
    connectionMutex.RLock()
    db, exists := connections[userID]
    connectionMutex.RUnlock()
    
    if exists && db.Ping() == nil {
        return db, nil
    }
    
    // Crear nueva conexión con lock exclusivo
    connectionMutex.Lock()
    defer connectionMutex.Unlock()
    
    // ... crear conexión
}
```

### 5. Timeouts y Límites

**Mejoras recomendadas**:

1. **Timeout de conexión**: 30 segundos
2. **Timeout de query**: 30 segundos por query
3. **Timeout total de auditoría**: 5 minutos
4. **Límite de conexiones simultáneas por usuario**: 1 (ya implementado con OneToOne)
5. **Pool de conexiones**: Máximo 25 conexiones abiertas, 5 idle

### 6. Logging y Monitoreo

**Recomendaciones para Golang**:

1. **Usar un logger estructurado**: `zap`, `logrus`, o `zerolog`
2. **Logging de operaciones críticas**:
   - Intentos de conexión (exitosos y fallidos)
   - Ejecución de queries
   - Errores de autenticación
   - Errores de base de datos

3. **Métricas**:
   - Tiempo de ejecución de queries
   - Número de conexiones activas
   - Tasa de error de conexiones
   - Tiempo de respuesta de endpoints

### 7. Validación de Datos

#### Validación de Entrada de Conexión

```go
type ConnectRequest struct {
    Driver   string `json:"driver" binding:"required,min=1,max=255"`
    Server   string `json:"server" binding:"required,min=1,max=255"`
    DBUser   string `json:"db_user" binding:"required,min=1,max=255"`
    Password string `json:"password" binding:"required,min=1"`
}

func (r *ConnectRequest) Validate() error {
    // Validar formato de IP o hostname (implementar isValidServer según necesidad)
    // Ejemplo básico:
    if r.Server == "" {
        return fmt.Errorf("server cannot be empty")
    }
    
    // Validar driver conocido
    validDrivers := []string{
        "ODBC Driver 17 for SQL Server",
        "ODBC Driver 13 for SQL Server",
        "SQL Server Native Client 11.0",
    }
    valid := false
    for _, d := range validDrivers {
        if d == r.Driver {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("invalid driver")
    }
    
    return nil
}
```

### 8. Manejo de Errores

#### Errores Comunes y Cómo Manejarlos

1. **Error de conexión a SQL Server**:
   - Credenciales incorrectas
   - Servidor inaccesible
   - Puerto bloqueado por firewall
   - Timeout de conexión

2. **Error de ejecución de query**:
   - Query SQL inválida
   - Permisos insuficientes
   - Tabla o columna no existe
   - Timeout de query

3. **Errores de autenticación**:
   - Token inválido o expirado
   - Usuario inactivo
   - Permisos insuficientes

**Estrategia de manejo en Golang**:
```go
import (
    "database/sql"
    "errors"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
)

type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func HandleError(c *gin.Context, err error) {
    var appErr *AppError
    
    errStr := err.Error()
    switch {
    case err == sql.ErrNoRows:
        appErr = &AppError{
            Code:    http.StatusNotFound,
            Message: "Resource not found",
        }
    case strings.Contains(errStr, "login failed"):
        appErr = &AppError{
            Code:    http.StatusUnauthorized,
            Message: "Invalid credentials",
        }
    case strings.Contains(errStr, "timeout"):
        appErr = &AppError{
            Code:    http.StatusRequestTimeout,
            Message: "Request timeout",
        }
    default:
        appErr = &AppError{
            Code:    http.StatusInternalServerError,
            Message: "Internal server error",
            Details: errStr, // Solo en desarrollo
        }
    }
    
    c.JSON(appErr.Code, appErr)
}
```

---

## Resumen de Endpoints Completos

### Tabla Resumen de Endpoints

| Método | Endpoint | Permisos | Descripción |
|--------|----------|----------|-------------|
| POST | `/api/users/register/` | AllowAny | Registrar nuevo usuario |
| POST | `/api/users/login/` | AllowAny | Iniciar sesión |
| GET | `/api/users/profile/` | IsAuthenticated | Obtener perfil |
| POST | `/api/users/logout/` | IsAuthenticated | Cerrar sesión |
| PUT | `/api/users/update_profile/` | IsAuthenticated | Actualizar perfil |
| POST | `/api/users/change_password/` | IsAuthenticated | Cambiar contraseña |
| POST | `/api/users/deactivate_account/` | IsAuthenticated | Desactivar cuenta |
| POST | `/api/sql_conn/connections/connect/` | IsAuthenticated, IsClient | Conectar a SQL Server |
| POST | `/api/sql_conn/connections/disconnect/` | IsAuthenticated, IsClient | Desconectar de SQL Server |
| GET | `/api/sql_conn/admin/active_connections/` | IsAdmin | Listar conexiones (admin) |
| GET | `/api/sql_conn/admin/active_connections/<id>/` | IsAdmin | Detalle de conexión (admin) |
| GET | `/api/sql/controls/controls_info/` | IsAuthenticated | Listar controles |
| GET | `/api/sql/controls/control_info/<id>/` | IsAuthenticated | Detalle de control |
| GET | `/api/sql/controls/execute/` | IsAuthenticated, IsClient | Ejecutar auditoría |
| GET | `/api/logs/connection_logs_list/` | IsClientAndOwner | Listar logs de conexión |
| GET | `/api/logs/connection_logs_details/<id>/` | IsClient | Detalle de log de conexión |
| GET | `/api/logs/admin/connection_logs_list/` | IsAdmin | Listar todos los logs (admin) |
| GET | `/api/logs/auditory_logs_list/` | IsClient | Listar auditorías |
| GET | `/api/logs/auditory_logs_detail/<id>` | IsClientAndOwner | Detalle de auditoría |
| GET | `/api/logs/auditory_logs_results/<audit_id>/` | IsClient | Resultados de auditoría |
| GET | `/api/logs/admin/auditory_logs_list/` | IsAdmin | Listar todas las auditorías (admin) |
| GET | `/api/dashGET/auditoryAmount/` | IsAuthenticated | Cantidad de auditorías |
| GET | `/api/dashGET/connectionAmount/` | IsAuthenticated | Cantidad de conexiones |
| GET | `/api/dashGET/correctRate/` | IsAuthenticated | Tasa de correctitud |

---

## Checklist de Migración a Golang

### Fase 1: Preparación
- [ ] Analizar todos los endpoints y su funcionalidad
- [ ] Mapear modelos de Django a structs de Go
- [ ] Diseñar estructura de proyecto Go
- [ ] Seleccionar framework web (Gin/Echo/Fiber)
- [ ] Configurar entorno de desarrollo Go

### Fase 2: Infraestructura Base
- [ ] Configurar base de datos SQLite con GORM
- [ ] Crear migraciones de base de datos
- [ ] Implementar sistema de autenticación JWT
- [ ] Configurar CORS
- [ ] Implementar middleware de autenticación
- [ ] Configurar logging estructurado

### Fase 3: Modelos y Servicios
- [ ] Implementar modelos (CustomUser, ActiveConnection, etc.)
- [ ] Implementar servicios de usuario
- [ ] Implementar servicios de conexión a SQL Server
- [ ] Implementar pool de conexiones
- [ ] Implementar servicios de auditoría
- [ ] Implementar servicios de logs

### Fase 4: Handlers/Controllers
- [ ] Implementar handlers de usuarios
- [ ] Implementar handlers de conexiones
- [ ] Implementar handlers de controles
- [ ] Implementar handlers de logs
- [ ] Implementar handlers de dashboard

### Fase 5: Funcionalidades Avanzadas
- [ ] Implementar encriptación de contraseñas SQL Server
- [ ] Implementar timeouts y límites
- [ ] Implementar validación de datos
- [ ] Implementar manejo de errores robusto
- [ ] Implementar métricas y monitoreo

### Fase 6: Testing
- [ ] Tests unitarios de modelos
- [ ] Tests unitarios de servicios
- [ ] Tests de integración de endpoints
- [ ] Tests de conexión a SQL Server
- [ ] Tests de ejecución de queries

### Fase 7: Despliegue
- [ ] Configurar variables de entorno
- [ ] Configurar Docker (opcional)
- [ ] Configurar CI/CD
- [ ] Documentación de API (Swagger/OpenAPI)
- [ ] Guía de despliegue

---

## Notas Finales

### Mejoras de Seguridad Recomendadas

1. **Encriptar contraseñas de SQL Server** antes de almacenar
2. **Usar JWT con refresh tokens** en lugar de tokens simples
3. **Implementar rate limiting** para prevenir ataques de fuerza bruta
4. **Validar y sanitizar todas las entradas** de usuario
5. **Usar prepared statements** para prevenir SQL injection
6. **Implementar HTTPS** en producción
7. **Configurar CORS** correctamente para producción
8. **Usar variables de entorno** para secretos y configuraciones
9. **Implementar logging de seguridad** (intentos de acceso, errores, etc.)
10. **Auditar regularmente** los logs de acceso

### Consideraciones de Rendimiento

1. **Pool de conexiones**: Configurar adecuadamente según carga esperada
2. **Caché de controles**: Los controles raramente cambian, considerar caché
3. **Ejecución paralela de queries**: Para auditorías completas, considerar ejecutar queries en paralelo (con límites)
4. **Índices de base de datos**: Asegurar índices en campos de búsqueda frecuente
5. **Paginación**: Implementar paginación para listados grandes

### Mantenimiento

1. **Versionado de API**: Considerar versionado de endpoints (`/api/v1/...`)
2. **Documentación**: Mantener documentación actualizada
3. **Monitoreo**: Implementar métricas y alertas
4. **Backups**: Configurar backups regulares de la base de datos
5. **Actualizaciones**: Mantener dependencias actualizadas

---

## Recursos Adicionales

### Documentación de Librerías Go

- **Gin**: https://gin-gonic.com/docs/
- **GORM**: https://gorm.io/docs/
- **go-mssqldb**: https://github.com/denisenkom/go-mssqldb
- **JWT**: https://github.com/golang-jwt/jwt
- **bcrypt**: https://pkg.go.dev/golang.org/x/crypto/bcrypt

### Referencias de SQL Server

- **Connection Strings**: https://www.connectionstrings.com/sql-server/
- **SQL Server Security**: https://docs.microsoft.com/en-us/sql/relational-databases/security/

---

**Documento generado para migración de Django (Python) a Golang**  
**Última actualización**: Enero 2025