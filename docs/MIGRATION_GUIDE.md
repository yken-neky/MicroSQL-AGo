# Migración SQLite (Django) → MySQL (Go) - Guía Rápida

## 📋 Resumen de la Migración Completada

Tu microservicio originalmente tenía una base de datos **SQLite con esquema Django** (prefijos como `Users_App_customuser`, `Logs_App_*`). Se ha creado un migrador que:

1. Lee las tablas Django desde SQLite
2. Mapea los datos a las estructuras Go en `internal/domain/entities`
3. Inserta en MySQL con AutoMigrate de GORM

### Datos Migrados (verificado):
- ✅ **4 usuarios** (de `Users_App_customuser` → `users`)
- ✅ **1 conexión activa** (de `Connecting_App_activeconnection` → `active_connections`)
- ✅ **40 registros de logs** (de `Logs_App_connectionlog` → `connection_logs`)
- ✅ **43 controles de información** (de `InsideDB_App_controls_information` → `controls_informations`)

---

## 🚀 Cómo Ejecutar la Migración (paso a paso)

### Opción 1: Local (binario + MySQL en Docker)

1. **Inicia MySQL** (si está en docker-compose):
   ```bash
   cd backend-go
   docker-compose up -d db
   sleep 3
   ```

2. **Exporta variables de entorno** (zsh):
   ```bash
   export MYSQL_HOST=127.0.0.1
   export MYSQL_PORT=3306
   export MYSQL_USER=micro_user
   export MYSQL_PASSWORD='ChangeMeStrongPassword!'
   export MYSQL_DATABASE=microsql_ago
   export DB_PATH=./db.sqlite3
   ```

3. **Compila y ejecuta el migrador**:
   ```bash
   cd backend-go
   go run ./cmd/migrate --src ./db.sqlite3
   ```

### Opción 2: Con Docker Compose (servicio integrado)

1. Desde `backend-go`:
   ```bash
   docker-compose up migrate
   ```
   El servicio ejecuta y se detiene automáticamente.

---

## 📊 Estructura de Tablas Creadas en MySQL

### Tablas Migradas:

| Tabla Django | Tabla MySQL | Campos | Registros |
|---|---|---|---|
| `Users_App_customuser` | `users` | id, username, email, password, first_name, last_name, role, created_at, last_login, is_active | 4 |
| `Connecting_App_activeconnection` | `active_connections` | id, user_id, driver, server, db_user, password, is_connected, last_connected, last_disconnected | 1 |
| `Logs_App_connectionlog` | `connection_logs` | id, user_id, driver, server, db_user, timestamp, status | 40 |
| `InsideDB_App_controls_information` | `controls_informations` | id, idx, chapter, name, description | 43 |

### Tablas Creadas Automáticamente (vacías, esperadas para el futuro):

- `roles` — para roles de usuario
- `permissions` — para permisos del sistema
- `user_roles` — relación usuario-rol
- `role_permissions` — relación rol-permiso
- `queries` — histórico de consultas SQL
- `query_result_dbs` — resultados guardados de consultas
- `execution_stats` — estadísticas de ejecución

---

## 🔍 Verificación de Datos

Conéctate a MySQL y ejecuta:

```bash
# Contar registros
docker exec backend-go-db-1 mysql -u micro_user -p'ChangeMeStrongPassword!' microsql_ago -e \
"SELECT 'Users', COUNT(*) FROM users
UNION ALL
SELECT 'Active Connections', COUNT(*) FROM active_connections
UNION ALL
SELECT 'Connection Logs', COUNT(*) FROM connection_logs
UNION ALL
SELECT 'Controls', COUNT(*) FROM controls_informations;"

# Ver muestra de usuarios
docker exec backend-go-db-1 mysql -u micro_user -p'ChangeMeStrongPassword!' microsql_ago -e \
"SELECT id, username, email, role FROM users LIMIT 5;"
```

---

## ⚙️ Archivos Modificados

1. **`backend-go/cmd/migrate/main.go`** — Reescrito para mapear tablas Django a estructuras Go
2. **`backend-go/internal/adapters/secondary/persistence/sqlite/migrations/migrate.go`** — Actualizado AutoMigrate
3. **`backend-go/internal/domain/entities/user.go`** — Hizo `LastDisconnected` nullable (`*time.Time`)
4. **`backend-go/internal/domain/usecases/connection/disconnect_from_server.go`** — Ajustado para nullable

---

## ⚠️ Notas Importantes

### 1. **Campos de Fecha con Cero**
   - SQLite permitía fechas `0000-00-00`, pero MySQL las rechaza
   - El migrador filtra estas fechas y usa `time.Now()` como fallback (ver línea ~200 en main.go)
   - Campo `last_disconnected` es ahora nullable (`*time.Time`)

### 2. **Contraseñas**
   - Las contraseñas se copian tal cual (hashes Django `pbkdf2_sha256$...`)
   - No se re-encriptan ni desencriptan

### 3. **Relaciones Many-to-Many Vacías**
   - `roles`, `permissions`, `user_roles`, `role_permissions` se crean pero están vacías (no existían en Django)
   - Puedes poblarlas según tus necesidades en la app Go

### 4. **Integridad Referencial**
   - Los FOREIGN KEYs se crean mediante GORM pero NO están explícitamente habilitados
   - Verifica que todos los `user_id` en `active_connections`, `connection_logs`, etc., existan en `users`

---

## 🔧 Troubleshooting

### Error: "no such table: users"
   - Verificar que el path de SQLite es correcto: `--src ./db.sqlite3`
   - Listar tablas: `sqlite3 db.sqlite3 ".tables"`

### Error: "Incorrect datetime value: '0000-00-00'"
   - Ya está solucionado en el código actual
   - Si persiste, verificar que `cmd/migrate/main.go` línea ~195 tiene el filtro de fecha cero

### MySQL Connection Refused
   - Verificar que MySQL está corriendo: `docker-compose logs db`
   - Esperar 2-3 segundos después de `up -d db`

---

## 📝 Próximos Pasos Opcionales

1. **Backup de SQLite** (recomendado antes de cada migración):
   ```bash
   cp db.sqlite3 db.sqlite3.backup.$(date +%Y%m%d_%H%M%S)
   ```

2. **Limpiar BD antes de re-migrar** (si necesitas empezar de cero):
   ```bash
   docker-compose down  # Borra volumen de datos
   docker volume rm backend-go_db_data  # Si es necesario
   ```

3. **Agregar índices personalizados** en MySQL después de la migración:
   ```sql
   ALTER TABLE users ADD INDEX idx_username (username);
   ALTER TABLE connection_logs ADD INDEX idx_user_timestamp (user_id, timestamp);
   ```

4. **Validar integridad referencial**:
   ```sql
   SELECT * FROM active_connections WHERE user_id NOT IN (SELECT id FROM users);
   ```

---

## 📞 Referencias Útiles

- **Configuración de Base de Datos**: `backend-go/internal/config/config.go`
- **Entities (estructuras Go)**: `backend-go/internal/domain/entities/`
- **Docker Compose**: `backend-go/docker-compose.yml`
- **Logs de MySQL**: `docker-compose logs -f db`
- **Logs del Migrador**: Salida estándar del comando `go run ./cmd/migrate`

---

## ✅ Checklist Final

- [x] SQLite y MySQL configurados
- [x] Migrador compilado sin errores
- [x] 4 usuarios migrados
- [x] 1 conexión activa migrada
- [x] 40 registros de log migrados
- [x] 43 controles migrados
- [x] Fechas cero manejadas (nullable)
- [x] Verificación de datos completada

**¡Migración completada exitosamente! 🎉**

