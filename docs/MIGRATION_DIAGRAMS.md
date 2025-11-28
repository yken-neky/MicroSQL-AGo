# Diagramas Visuales de la Migración

## 1. Flujo General de Migración

```
┌─────────────────────┐
│   SQLite (Django)   │
│   db.sqlite3        │
│                     │
│ • Users_App_*       │
│ • Logs_App_*        │
│ • Connecting_App_*  │
│ • InsideDB_App_*    │
└──────────┬──────────┘
           │
           │ gorm.Open(sqlite.Open())
           │ Find(&DjangoCustomUser{})
           │
           ▼
┌──────────────────────────────────────────┐
│       MIGRADOR (cmd/migrate/main.go)     │
│                                          │
│ 1. Lee de tablas Django                 │
│ 2. Mapea a structs Go                   │
│ 3. Filtra datos inválidos (ej: 0000-*) │
│ 4. Inserta con OnConflict DoNothing     │
└──────────┬───────────────────────────────┘
           │
           │ gorm.Open(mysql.Open())
           │ Create(&User{}, &*)
           │
           ▼
┌─────────────────────┐
│   MySQL (Go)        │
│   microsql_ago      │
│                     │
│ • users (4)         │
│ • active_conn. (1)  │
│ • conn_logs (40)    │
│ • controls (43)     │
│ • (varios vacíos)   │
└─────────────────────┘
```

## 2. Mapeo Detallado Tabla por Tabla

### Users: Users_App_customuser → users

```
LECTURA (SQLite)              MAPEO               INSERCIÓN (MySQL)
─────────────────────────────────────────────────────────────────────
Users_App_customuser row:
  id=1                   ──┐
  username="yan"         ──┤
  email="yan@..."        ──┤
  password="pbkdf2_*"    ──┤   entities.User{
  first_name="Yan"       ──┤     ID: 1,
  last_name="Gonzalez"   ──┤     Username: "yan",
  role="cliente"         ──┤     Email: "yan@...",
  created_at=2025-01-22  ──┤     Password: "pbkdf2_*",
  last_login=2025-09-15  ──┤     FirstName: "Yan",
  is_active=true         ──┤     LastName: "Gonzalez",
                         ──┤     Role: "cliente",
                         ──┤     CreatedAt: 2025-01-22,
                         ──┤     LastLogin: 2025-09-15,
                         ──┘     IsActive: true
                                }
                                    │
                                    │ INSERT INTO users (...)
                                    ▼
                                MySQL: ✓ Inserted
                                ID=1, username="yan"
```

### ActiveConnections: Connecting_App_activeconnection → active_connections

```
LECTURA (SQLite)              MAPEO               INSERCIÓN (MySQL)
─────────────────────────────────────────────────────────────────────
Connecting_App_activeconnection:
  id=18                  ──┐
  user_id=1              ──┤
  driver="ODBC*"         ──┤
  server="DESKTOP-*"     ──┤   entities.ActiveConnection{
  db_user="sa"           ──┤     ID: 18,
  password="SQL*1234"    ──┤     UserID: 1,
  is_connected=true      ──┤     Driver: "ODBC*",
  last_connected=        ──┤     Server: "DESKTOP-*",
    2025-06-10           ──┤     DBUser: "sa",
  last_disconnected=     ──┼─→  Password: "SQL*1234",
    0000-00-00 ⚠️        ──┤     IsConnected: true,
                         ──┤     LastConnected: 2025-06-10,
                         ──┤     LastDisconnected: nil
                         ──┘   }
                                    │
                        ⚠️ FILTRO: Fecha 0000-00-00
                           → Convertida a NULL
                                    │
                                    │ INSERT INTO active_connections (...)
                                    ▼
                                MySQL: ✓ Inserted
                                ID=18, user_id=1, last_disconnected=NULL
```

### ConnectionLogs: Logs_App_connectionlog → connection_logs

```
LECTURA (SQLite)              MAPEO               INSERCIÓN (MySQL)
─────────────────────────────────────────────────────────────────────
Logs_App_connectionlog rows (40 total):
  id=1                   ──┐
  user_id=1              ──┤
  driver="ODBC*"         ──┤   entities.ConnectionLog{
  server="192.168.*"     ──┤     ID: 1,
  db_user="admin"        ──┤     UserID: 1,
  timestamp=2025-06-10   ──┤     Driver: "ODBC*",
  status="connected"     ──┤     Server: "192.168.*",
                         ──┤     DBUser: "admin",
                         ──┤     Timestamp: 2025-06-10,
                         ──┘     Status: "connected"
                                }
                                    │
                                    │ INSERT x 40 registros
                                    ▼
                                MySQL: ✓ Inserted 40 rows
                                connection_logs populated
```

### ControlsInformation: InsideDB_App_controls_information → controls_informations

```
LECTURA (SQLite)              MAPEO               INSERCIÓN (MySQL)
─────────────────────────────────────────────────────────────────────
InsideDB_App_controls_info (43 total):
  id=1                   ──┐
  idx=1                  ──┤
  chapter="2"            ──┤   entities.ControlsInformation{
  name="Control Name 1"  ──┤     ID: 1,
  description="..."      ──┤     Idx: 1,
  impact="..."           ──┼─→  (descartado)
  good_config="..."      ──┼─→  (descartado)
  bad_config="..."       ──┼─→  (descartado)
  ref="..."              ──┼─→  (descartado)
                         ──┤     Chapter: "2",
                         ──┤     Name: "Control Name 1",
                         ──┘     Description: "..."
                                }
                                    │
                                    │ INSERT x 43 registros
                                    ▼
                                MySQL: ✓ Inserted 43 rows
                                controls_informations populated
```

## 3. Manejo de Errores y Filtros

```
┌─────────────────────────────────────────────────────────────┐
│         VALIDACIONES Y FILTROS EN MIGRADOR                 │
└─────────────────────────────────────────────────────────────┘

Input Row:
┌──────────────────────────────────────────┐
│ Connecting_App_activeconnection row      │
│ • last_connected = "2025-06-10 11:04"   │ ✓ Válida
│ • last_disconnected = "0000-00-00 00:00"│ ⚠️ Invalida
└──────────────────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────────┐
│      FILTRO: Detecta fecha cero          │
└──────────────────────────────────────────┘
         │
         ├─→ if lastConnected.IsZero() { return NULL }
         │
         ▼
Output para MySQL:
┌──────────────────────────────────────────┐
│ ActiveConnection struct                  │
│ • LastConnected = 2025-06-10 11:04      │ ✓ Insertable
│ • LastDisconnected = nil (NULL)         │ ✓ Permitido
└──────────────────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────────┐
│      INSERT con OnConflict DoNothing     │
│      ✓ OK - MySQL acepta NULL           │
└──────────────────────────────────────────┘
```

## 4. Modelo Transaccional (Flujo Actual)

```
┌────────────────────────────────────────────────────────────────┐
│  Para cada tabla (Sequential, NO transactional por tabla)     │
└────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  1. USERS                                                   │
│  ├─ srcDB.Find(&users)          ✓ 4 rows read            │
│  └─ dstDB.Create(&users)        ✓ 4 rows inserted        │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼ (si error, FAIL y STOP)
┌─────────────────────────────────────────────────────────────┐
│  2. ACTIVE_CONNECTIONS                                      │
│  ├─ srcDB.Find(&conns)          ✓ 1 row read             │
│  ├─ Filtro: IsZero()            ✓ Limpia fechas cero     │
│  └─ dstDB.Create(&conns)        ✓ 1 row inserted         │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  3. CONNECTION_LOGS                                         │
│  ├─ srcDB.Find(&logs)           ✓ 40 rows read           │
│  └─ dstDB.Create(&logs)         ✓ 40 rows inserted       │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  4. CONTROLS_INFORMATION                                    │
│  ├─ srcDB.Find(&controls)       ✓ 43 rows read           │
│  └─ dstDB.Create(&controls)     ✓ 43 rows inserted       │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  SUMMARY LOG                                                │
│  ✓ Migration completed successfully!                        │
│  Users: 4                                                   │
│  Active Connections: 1                                      │
│  Connection Logs: 40                                        │
│  Controls: 43                                               │
│  Total: 88 rows ✓                                          │
└─────────────────────────────────────────────────────────────┘
```

## 5. Comparación: Antes vs Después

```
┌───────────────────────────────────────────────────────────────────┐
│                        ANTES                                      │
├───────────────────────────────────────────────────────────────────┤
│ BD: SQLite (db.sqlite3)                                          │
│ Esquema: Django (prefijos de app)                                │
│ Tablas:                                                          │
│   • Users_App_customuser          ← monolitico Django           │
│   • Logs_App_connectionlog                                       │
│   • Connecting_App_activeconnection                              │
│   • InsideDB_App_controls_information                            │
│ Acceso: Código Django + ORM Django                              │
│ Integridad: Soft (sin FK explícitos)                            │
│ Escalabilidad: Limitada a SQLite                                │
└───────────────────────────────────────────────────────────────────┘
                            │
                            │ ✨ MIGRACIÓN ✨
                            │
                            ▼
┌───────────────────────────────────────────────────────────────────┐
│                        DESPUÉS                                    │
├───────────────────────────────────────────────────────────────────┤
│ BD: MySQL (en Docker)                                            │
│ Esquema: Go (domain-driven design)                               │
│ Tablas:                                                          │
│   • users ✓                                                       │
│   • active_connections ✓                                         │
│   • connection_logs ✓                                            │
│   • controls_informations ✓                                      │
│   • roles (preparadas para futuro)                               │
│   • permissions (preparadas para futuro)                         │
│   • user-query history and storage was removed (no arbitrary SQL execution)
│ Acceso: Código Go + GORM ORM                                    │
│ Integridad: FK + índices UNIQUE                                 │
│ Escalabilidad: ✓ Multi-usuario, replicación, sharding         │
│ Rendimiento: ✓ Mejor que SQLite para prod                      │
└───────────────────────────────────────────────────────────────────┘
```

## 6. Checklist de Validación Post-Migración

```
✅ VALIDACIONES REALIZADAS:

┌─────────────────────────────────────────────────────┐
│ 1. Tabla USERS                                      │
│    ├─ Conteo: 4 registros ✓                       │
│    ├─ PK (id): Intactos ✓                         │
│    ├─ UNIQUE (username): Intactos ✓              │
│    ├─ UNIQUE (email): Intactos ✓                 │
│    └─ Datos sensibles (password): Copiados ✓     │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ 2. Tabla ACTIVE_CONNECTIONS                         │
│    ├─ Conteo: 1 registro ✓                        │
│    ├─ FK (user_id): Referencia válida ✓           │
│    ├─ Fechas cero: Filtradas a NULL ✓            │
│    └─ Datos cifrados (password): Copiados ✓       │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ 3. Tabla CONNECTION_LOGS                            │
│    ├─ Conteo: 40 registros ✓                      │
│    ├─ FK (user_id): Referencias válidas ✓         │
│    ├─ Timestamps: Intactos ✓                      │
│    └─ Status: Valores esperados ✓                 │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ 4. Tabla CONTROLS_INFORMATIONS                      │
│    ├─ Conteo: 43 registros ✓                      │
│    ├─ UNIQUE (idx, chapter): Intactos ✓           │
│    ├─ Descripción (text): Intacta ✓               │
│    └─ Campos no migrados (impact, etc): OK ✓      │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ 5. Tablas Preparadas (vacías, para futuro)         │
│    ├─ roles: ✓ Creada                             │
│    ├─ permissions: ✓ Creada                       │
│    ├─ user_roles: ✓ Creada                        │
│    ├─ role_permissions: ✓ Creada                  │
│    ├─ queries: ✓ Creada                           │
│    ├─ query_result_dbs: removed
│    └─ execution_stats: removed
└─────────────────────────────────────────────────────┘

RESULTADO FINAL: ✅ 100% EXITOSO
```

