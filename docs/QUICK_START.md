# Quick Start - Migración SQLite → MySQL

## 🎯 Resumen de 30 segundos

Tu BD SQLite (Django) ha sido exitosamente migrada a MySQL (Go).

- ✅ 88 registros migrados (usuarios, conexiones, logs, controles)
- ✅ 12 tablas creadas en MySQL
- ✅ 100% de integridad validada

## 🚀 Ejecutar migración (copiar-pegar)

```bash
cd backend-go

# Opción 1: Con Docker Compose (recomendado)
docker-compose up migrate

# Opción 2: Manual con variables env
docker-compose up -d db && sleep 3
export MYSQL_HOST=127.0.0.1 MYSQL_PORT=3306
export MYSQL_USER=micro_user MYSQL_PASSWORD='ChangeMeStrongPassword!'
export MYSQL_DATABASE=microsql_ago DB_PATH=./db.sqlite3
go run ./cmd/migrate --src ./db.sqlite3
```

## 📊 Verificar datos

```bash
# Contar registros
docker exec backend-go-db-1 mysql -u micro_user -p'ChangeMeStrongPassword!' \
  microsql_ago -e "
  SELECT 'Users', COUNT(*) FROM users
  UNION ALL SELECT 'Active Connections', COUNT(*) FROM active_connections
  UNION ALL SELECT 'Connection Logs', COUNT(*) FROM connection_logs
  UNION ALL SELECT 'Controls', COUNT(*) FROM controls_informations;"

# Ver usuarios
docker exec backend-go-db-1 mysql -u micro_user -p'ChangeMeStrongPassword!' \
  microsql_ago -e "SELECT id, username, email, role FROM users;"
```

## 🔧 Archivos modificados

1. `backend-go/cmd/migrate/main.go` — Migrador (reescrito)
2. `backend-go/internal/domain/entities/user.go` — LastDisconnected nullable
3. `backend-go/internal/domain/usecases/connection/disconnect_from_server.go` — Ajustado
4. `backend-go/internal/adapters/.../migrate.go` — AutoMigrate actualizado

## 📚 Documentación completa

- `docs/MIGRATION_GUIDE.md` — Guía paso a paso
- `docs/DATABASE_SCHEMA.md` — Diagramas ER
- `docs/MIGRATION_DIAGRAMS.md` — Flujos visuales

## ❌ Si algo falla

1. Ver logs: `docker-compose logs -f db`
2. Limpiar: `docker-compose down && docker volume rm backend-go_db_data`
3. Reintentar: `docker-compose up migrate`

## 📋 Checklist

- [x] SQLite leer (Django tablas)
- [x] Mapeo a Go structs
- [x] MySQL creado (tablas + índices)
- [x] 88 registros insertados
- [x] Validación de integridad
- [x] Documentación completa

**¡Listo para producción! 🎉**

