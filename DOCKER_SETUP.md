# 🐳 Configuración Docker - MicroSQL AGo

Este documento explica cómo ejecutar toda la aplicación usando Docker Compose.

## 📋 Requisitos Previos

- Docker Engine 20.10+
- Docker Compose 2.0+

## 🚀 Inicio Rápido

### 1. Construir y levantar todos los servicios

```bash
docker-compose up -d --build
```

Este comando:
- Construye las imágenes del backend y frontend
- Levanta MySQL, backend y frontend
- Configura las redes y volúmenes necesarios

### 2. Ejecutar migraciones (primera vez)

```bash
docker-compose --profile migration run --rm migrate
```

### 3. Acceder a la aplicación

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8000
- **MySQL:** localhost:3306

## 📦 Servicios Incluidos

### 1. **db** (MySQL 8.0)
- Base de datos principal
- Puerto: `3306`
- Usuario: `micro_user`
- Contraseña: `ChangeMeStrongPassword!`
- Base de datos: `microsql_ago`

### 2. **backend** (Go API)
- API REST en Go
- Puerto: `8000`
- Espera a que MySQL esté saludable antes de iniciar
- Variables de entorno configurables

### 3. **frontend** (Next.js)
- Aplicación web en React/Next.js
- Puerto: `3000`
- Se conecta al backend en `http://localhost:8000`

### 4. **migrate** (Migraciones)
- Ejecuta migraciones de base de datos
- Solo se ejecuta cuando se usa el perfil `migration`
- Se ejecuta una vez y termina

## 🔧 Comandos Útiles

### Ver logs de todos los servicios
```bash
docker-compose logs -f
```

### Ver logs de un servicio específico
```bash
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f db
```

### Detener todos los servicios
```bash
docker-compose down
```

### Detener y eliminar volúmenes (⚠️ elimina datos)
```bash
docker-compose down -v
```

### Reconstruir un servicio específico
```bash
docker-compose build backend
docker-compose up -d backend
```

### Ejecutar comandos dentro de un contenedor
```bash
# Backend
docker-compose exec backend sh

# Frontend
docker-compose exec frontend sh

# MySQL
docker-compose exec db mysql -u micro_user -p microsql_ago
```

## 🔐 Variables de Entorno

### Backend

Las siguientes variables pueden ser configuradas en el archivo `.env` o directamente en `docker-compose.yml`:

```env
# Server
SERVER_PORT=8000
GIN_MODE=release
LOG_LEVEL=info

# MySQL (ya configuradas en docker-compose.yml)
MYSQL_HOST=db
MYSQL_PORT=3306
MYSQL_USER=micro_user
MYSQL_PASSWORD=ChangeMeStrongPassword!
MYSQL_DATABASE=microsql_ago

# Security (⚠️ Cambiar en producción!)
JWT_SECRET=change-me-in-production-use-strong-secret
ENCRYPTION_KEY=01234567890123456789012345678901

# SQLite fallback
DB_PATH=/app/db.sqlite3
```

### Frontend

```env
NEXT_PUBLIC_API_URL=http://localhost:8000
NODE_ENV=production
PORT=3000
```

## 🛠️ Desarrollo

### Modo desarrollo con hot-reload

Para desarrollo, es recomendable ejecutar los servicios localmente:

**Backend:**
```bash
cd backend-go
go run cmd/server/main.go
```

**Frontend:**
```bash
cd frontend-nextjs
npm run dev
```

Y solo usar Docker para MySQL:
```bash
docker-compose up db
```

### Reconstruir después de cambios

Si haces cambios en el código y quieres reconstruir:

```bash
# Reconstruir todo
docker-compose up -d --build

# O solo un servicio
docker-compose build frontend
docker-compose up -d frontend
```

## 🐛 Troubleshooting

### El backend no puede conectarse a MySQL

1. Verifica que MySQL esté saludable:
   ```bash
   docker-compose ps
   ```

2. Revisa los logs de MySQL:
   ```bash
   docker-compose logs db
   ```

3. Espera unos segundos después de `docker-compose up` para que MySQL termine de inicializarse

### El frontend no puede conectarse al backend

1. Verifica que el backend esté corriendo:
   ```bash
   curl http://localhost:8000/health
   ```

2. Revisa los logs del backend:
   ```bash
   docker-compose logs backend
   ```

3. Verifica que `NEXT_PUBLIC_API_URL` esté configurado correctamente

### Puerto ya en uso

Si obtienes un error de puerto en uso:

```bash
# Ver qué está usando el puerto
lsof -i :3000
lsof -i :8000
lsof -i :3306

# Cambiar los puertos en docker-compose.yml
```

### Limpiar todo y empezar de nuevo

```bash
# Detener y eliminar contenedores, redes y volúmenes
docker-compose down -v

# Eliminar imágenes también
docker-compose down -v --rmi all

# Limpiar sistema Docker (⚠️ elimina todo lo no usado)
docker system prune -a --volumes
```

## 📝 Notas

- Los datos de MySQL se persisten en el volumen `db_data`
- El archivo SQLite del backend se monta desde `./backend-go/db.sqlite3`
- El frontend usa el modo `standalone` de Next.js para optimizar el tamaño de la imagen
- Las migraciones se ejecutan manualmente con el perfil `migration`

## 🔄 Actualizar la aplicación

```bash
# 1. Detener servicios
docker-compose down

# 2. Actualizar código (git pull, etc.)

# 3. Reconstruir y levantar
docker-compose up -d --build

# 4. Ejecutar migraciones si hay cambios en BD
docker-compose --profile migration run --rm migrate
```

