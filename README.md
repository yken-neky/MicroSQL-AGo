# MicroSQL AGo

Plataforma completa para gestionar conexiones a bases de datos y ejecutar auditorías de seguridad.

## 📁 Estructura del Proyecto

```
MicroSQL-AGo/
├── backend-go/          # Backend en Go (API REST)
│   ├── cmd/            # Aplicaciones ejecutables
│   ├── internal/       # Código interno del backend
│   └── Dockerfile      # Imagen Docker del backend
│
├── frontend-nextjs/    # Frontend en Next.js (React/TypeScript)
│   ├── src/           # Código fuente del frontend
│   └── Dockerfile     # Imagen Docker del frontend
│
├── docs/              # Documentación del proyecto
├── docker-compose.yml # Orquestación de servicios Docker
└── README.md          # Este archivo
```

## 🚀 Inicio Rápido con Docker

La forma más fácil de ejecutar toda la aplicación es usando Docker Compose:

```bash
# 1. Construir y levantar todos los servicios
docker-compose up -d --build

# 2. Ejecutar migraciones (primera vez)
docker-compose --profile migration run --rm migrate

# 3. Acceder a la aplicación
# Frontend: http://localhost:3000
# Backend API: http://localhost:8000
```

Para más detalles sobre Docker, consulta [DOCKER_SETUP.md](./DOCKER_SETUP.md).

## 🛠️ Desarrollo Local

### Backend (Go)

```bash
cd backend-go

# Instalar dependencias
go mod download

# Ejecutar servidor
go run cmd/server/main.go

# El servidor estará en http://localhost:8000
```

### Frontend (Next.js)

```bash
cd frontend-nextjs

# Instalar dependencias
npm install

# Ejecutar servidor de desarrollo
npm run dev

# La aplicación estará en http://localhost:3000
```

### Base de Datos

Puedes usar MySQL con Docker:

```bash
docker-compose up db
```

O usar SQLite (por defecto en desarrollo).

## 📚 Documentación

- [DOCKER_SETUP.md](./DOCKER_SETUP.md) - Guía completa de Docker
- [ENDPOINTS_SUMMARY.md](./ENDPOINTS_SUMMARY.md) - Resumen de todos los endpoints de la API
- [docs/](./docs/) - Documentación adicional del proyecto

## 🔧 Tecnologías

### Backend
- **Go 1.24+** - Lenguaje de programación
- **Gin** - Framework web
- **GORM** - ORM para bases de datos
- **MySQL/SQLite** - Bases de datos soportadas
- **JWT** - Autenticación

### Frontend
- **Next.js 16** - Framework React
- **React 19** - Biblioteca de UI
- **TypeScript** - Tipado estático
- **Tailwind CSS 4** - Estilos

## 📝 Endpoints Principales

### Autenticación
- `POST /api/auth/login` - Iniciar sesión
- `POST /api/auth/logout` - Cerrar sesión
- `POST /api/users/register` - Registro de usuarios

### Conexiones a BD
- `POST /api/db/:manager/open` - Abrir conexión
- `DELETE /api/db/:manager/close` - Cerrar conexión
- `GET /api/db/connections` - Listar conexiones activas

### Auditorías
- `POST /api/db/:manager/audits/execute` - Ejecutar auditoría
- `GET /api/db/:manager/audits/:id` - Ver resultados

### Administración (requiere rol admin)
- `GET /api/admin/users` - Listar usuarios
- `GET /api/admin/roles` - Gestión de roles
- `GET /api/admin/metrics/*` - Métricas del sistema

Ver [ENDPOINTS_SUMMARY.md](./ENDPOINTS_SUMMARY.md) para la lista completa.

## 🔐 Seguridad

- Autenticación basada en JWT
- Contraseñas hasheadas con bcrypt
- Encriptación de credenciales de BD con AES-GCM
- Control de acceso basado en roles (RBAC)
- Política de sesión única

## 📦 Servicios Docker

- **db** - MySQL 8.0 (puerto 3306)
- **backend** - API Go (puerto 8000)
- **frontend** - Next.js (puerto 3000)
- **migrate** - Ejecutor de migraciones (perfil migration)

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Ver [LICENSE](./LICENSE) para más detalles.

## 🆘 Soporte

Para problemas o preguntas:
1. Revisa la documentación en `docs/`
2. Consulta [DOCKER_SETUP.md](./DOCKER_SETUP.md) para problemas con Docker
3. Revisa los logs: `docker-compose logs -f`

