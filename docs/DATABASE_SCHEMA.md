# Diagrama ER - Estructura de BD (Post-Migración)

## Vista General de Tablas

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          BASE DE DATOS: microsql_ago                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────┐
│            USERS (4 registros)           │
├──────────────────────────────────────────┤
│ id (PK)                    uint          │
│ username (UNIQUE)          varchar(150)  │
│ email (UNIQUE)             varchar(254)  │
│ password                   varchar(255)  │
│ first_name                 varchar(150)  │
│ last_name                  varchar(150)  │
│ role                       varchar(20)   │
│ created_at                 datetime      │
│ last_login                 datetime      │
│ is_active                  bool          │
└──────────────────────────────────────────┘
         ▲            ▲            ▲
         │            │            │
         ├────────────┼────────────┤
         │            │            │
         1 : N        1 : N        1 : 1
         │            │            │
         ▼            ▼            ▼

┌─────────────────────────┐    ┌────────────────────────────┐    ┌──────────────────────────┐
│  CONNECTION_LOGS        │    │  ACTIVE_CONNECTIONS       │    │  USER_ROLES              │
│ (40 registros)          │    │  (1 registro)             │    │  (vacía - relación)      │
├─────────────────────────┤    ├────────────────────────────┤    ├──────────────────────────┤
│ id (PK)                 │    │ id (PK)                    │    │ user_id (PK)             │
│ user_id (FK)            │    │ user_id (FK, UNIQUE)       │    │ role_id (PK)             │
│ driver    varchar(255)  │    │ driver        varchar(255) │    └──────────────────────────┘
│ server    varchar(255)  │    │ server        varchar(255) │
│ db_user   varchar(255)  │    │ db_user       varchar(255) │
│ timestamp datetime      │    │ password      varchar(255) │
│ status    varchar(50)   │    │ is_connected  bool         │
└─────────────────────────┘    │ last_connected   datetime  │
                               │ last_disconnected *datetime│
                               │ (nullable)                 │
                               └────────────────────────────┘

┌──────────────────────────────────────┐
│     CONTROLS_INFORMATIONS            │
│     (43 registros)                   │
├──────────────────────────────────────┤
│ id (PK)                   uint       │
│ idx (UNIQUE con chapter)  int        │
│ chapter (UNIQUE con idx)  varchar10  │
│ name                      varchar255 │
│ description               text       │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│           ROLES (vacía)              │
├──────────────────────────────────────┤
│ id (PK)                   uint       │
│ name (UNIQUE)             varchar    │
│ description               text       │
└──────────────────────────────────────┘
         │
         │ N : M (via role_permissions)
         │
         ▼

┌──────────────────────────────────────┐
│       PERMISSIONS (vacía)            │
├──────────────────────────────────────┤
│ id (PK)                   uint       │
│ name (UNIQUE)             varchar    │
│ description               text       │
│ resource                  varchar    │
│ action                    varchar    │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│     ROLE_PERMISSIONS (vacía)         │
├──────────────────────────────────────┤
│ role_id (FK)              uint       │
│ permission_id (FK)        uint       │
│ (Composite PK)                       │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│        QUERIES (vacía)               │
├──────────────────────────────────────┤
│ id (PK)                   uint       │
│ user_id (FK)              uint       │
│ connection_id (FK)        uint       │
│ sql                       text       │
│ status                    varchar    │
│ start_time                datetime   │
│ end_time                  datetime   │
│ rows_affected             int64      │
│ error                     text       │
│ database                  varchar    │
└──────────────────────────────────────┘
         │
         │ 1 : 1
         │
         ▼

┌──────────────────────────────────────┐
│    QUERY_RESULT_DBS (vacía)          │
├──────────────────────────────────────┤
│ query_id (FK)             uint       │
│ columns (JSON)            text       │
│ types (JSON)              text       │
│ rows (JSON)               text       │
│ has_more_rows             bool       │
│ page_size                 int        │
│ page_number               int        │
│ created_at                datetime   │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│    EXECUTION_STATS (vacía)           │
├──────────────────────────────────────┤
│ query_id (FK)             uint       │
│ duration_ms               float64    │
│ rows_affected             int64      │
│ cpu_time_ms               float64    │
│ io_time_ms                float64    │
│ memory_kb                 int64      │
└──────────────────────────────────────┘

```

---

## Mapeo Django → Go

```
Django (SQLite)                          Go (MySQL)
─────────────────────────────────────────────────────────────

Users_App_customuser                     users
├─ id                                    ├─ id
├─ username (unique)                     ├─ username (unique)
├─ email (unique)                        ├─ email (unique)
├─ password                              ├─ password
├─ first_name                            ├─ first_name
├─ last_name                             ├─ last_name
├─ role (default: "cliente")             ├─ role
├─ created_at                            ├─ created_at
├─ last_login (nullable)                 ├─ last_login
└─ is_active (default: true)             └─ is_active

Connecting_App_activeconnection          active_connections
├─ id                                    ├─ id
├─ user_id (FK, unique)                  ├─ user_id (FK, unique)
├─ driver                                ├─ driver
├─ server                                ├─ server
├─ db_user                               ├─ db_user
├─ password                              ├─ password
├─ is_connected                          ├─ is_connected
└─ last_connected                        ├─ last_connected
                                         └─ last_disconnected (NULL si no disponible)

Logs_App_connectionlog                   connection_logs
├─ id                                    ├─ id
├─ user_id (FK)                          ├─ user_id (FK)
├─ driver                                ├─ driver
├─ server                                ├─ server
├─ db_user                               ├─ db_user
├─ timestamp                             ├─ timestamp
└─ status                                └─ status

InsideDB_App_controls_information        controls_informations
├─ id                                    ├─ id
├─ idx                                   ├─ idx
├─ chapter                               ├─ chapter
├─ name                                  ├─ name
├─ description                           ├─ description
├─ impact                                (descartado en migración)
├─ good_config                           (descartado en migración)
├─ bad_config                            (descartado en migración)
└─ ref                                   (descartado en migración)

(Django no tenía)                        roles
                                         ├─ id
                                         ├─ name
                                         └─ description

(Django no tenía)                        permissions
                                         ├─ id
                                         ├─ name
                                         ├─ description
                                         ├─ resource
                                         └─ action

```

---

## Flujo de Migración en cmd/migrate

```
┌─────────────────────────────────────────────────────────────────┐
│         MIGRADOR: cmd/migrate/main.go                           │
└─────────────────────────────────────────────────────────────────┘

1. Lee configuración
   ├─ DB_PATH (SQLite source)
   ├─ MYSQL_* env vars (MySQL dest)
   └─ Flags: --src, --dst-dsn

2. Abre BD SQLite (origen con esquema Django)
   └─ Struct mappers: DjangoCustomUser, DjangoConnectionLog, etc.

3. Abre BD MySQL (destino)
   └─ Crea conexión usando credenciales

4. AutoMigrate: Crea todas las tablas en MySQL
   ├─ entities.User, entities.ActiveConnection, etc.
   └─ repositories.QueryResultDB

5. Migra tablas EN ORDEN:
   ├─ Users:
   │   └─ Lee Users_App_customuser
   │   └─ Mapea a entities.User
   │   └─ Inserta en MySQL (OnConflict DoNothing)
   │   └─ ✓ 4 rows
   │
   ├─ Active Connections:
   │   └─ Lee Connecting_App_activeconnection
   │   └─ Filtra fechas nulas (MySQL rechaza 0000-00-00)
   │   └─ Mapea a entities.ActiveConnection
   │   └─ Inserta en MySQL
   │   └─ ✓ 1 row
   │
   ├─ Connection Logs:
   │   └─ Lee Logs_App_connectionlog
   │   └─ Mapea a entities.ConnectionLog
   │   └─ Inserta en MySQL
   │   └─ ✓ 40 rows
   │
   └─ Controls Information:
       └─ Lee InsideDB_App_controls_information
       └─ Mapea a entities.ControlsInformation
       └─ Inserta en MySQL
       └─ ✓ 43 rows

6. Resumen e imprime totales
```

---

## Índices Automáticos (por GORM)

```
Tabla                     Índice                      Tipo
────────────────────────────────────────────────────────────
users                     idx_username                UNIQUE
users                     idx_email                   UNIQUE
users                     idx_role                    INDEX
users                     idx_is_active               INDEX

active_connections        idx_user_id                 UNIQUE
active_connections        idx_is_connected            INDEX

connection_logs           idx_user_id                 INDEX
connection_logs           idx_timestamp               INDEX

controls_informations     idx_idx_chapter             UNIQUE (composite)
controls_informations     idx_chapter                 INDEX
```

---

## Consideraciones de Integridad Referencial

### ✅ Validated (presente en datos)
- Todos los `user_id` en `active_connections` existen en `users`
- Todos los `user_id` en `connection_logs` existen en `users`

### ⚠️ Empty but Expected (creadas para futuro uso)
- `roles` vacía
- `permissions` vacía
- `user_roles` vacía (relación N:M usuarios-roles)
- `role_permissions` vacía (relación N:M roles-permisos)
- `queries` vacía
- `query_result_dbs` vacía
- `execution_stats` vacía

### 🔍 Foreign Key Soft Enforcement
- GORM crea columnas FK pero NO habilita constraints explícitos en MySQL por defecto
- Para habilitar checks estrictos (opcional), ejecuta post-migración:
  ```sql
  ALTER TABLE active_connections 
    ADD CONSTRAINT fk_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
  
  ALTER TABLE connection_logs 
    ADD CONSTRAINT fk_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
  ```

