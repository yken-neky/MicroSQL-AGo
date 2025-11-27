# Manual de Usuario

## Índice
1. [Introducción](#introducción)
2. [Primeros Pasos](#primeros-pasos)
3. [Gestión de Conexiones](#gestión-de-conexiones)
4. [Ejecución de Consultas](#ejecución-de-consultas)
5. [Gestión de Usuarios](#gestión-de-usuarios)
6. [Roles y Permisos](#roles-y-permisos)
7. [Troubleshooting](#troubleshooting)

## Introducción

MicroSQL AGo es una aplicación que permite gestionar conexiones y ejecutar consultas en SQL Server de manera segura y eficiente. Este manual te guiará a través de todas las funcionalidades disponibles.

## Primeros Pasos

### Login
1. Acceder a la API: `http://your-server:8080`
2. Autenticarse:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@ejemplo.com",
    "password": "contraseña"
  }'
```
3. Guardar el token JWT recibido

### Configuración Inicial
- Verificar permisos asignados
- Configurar preferencias de usuario
- Revisar límites de conexiones

## Gestión de Conexiones

### Crear Nueva Conexión
```bash
curl -X POST http://localhost:8080/api/db/mssql/open \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "server": "localhost",
    "port": "1433",
    "db_user": "sa",
    "password": "YourStrong!Passw0rd"
  }'
```

### Listar Conexiones Activas
```bash
curl -X GET http://localhost:8080/api/db/connections \
  -H "Authorization: Bearer <token>"
```

### Desconectar
```bash
curl -X DELETE http://localhost:8080/api/db/mssql/close \
  -H "Authorization: Bearer <token>"
```

## Ejecución de Consultas

### Ejecutar Consulta Simple
```bash
curl -X POST http://localhost:8080/api/queries/execute \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "SELECT * FROM users",
    "database": "master",
    "pageSize": 100
  }'
```

### Resultados Paginados
- Usar parámetros `pageSize` y `page`
- Verificar `hasMoreRows` en la respuesta
- Siguiente página: incrementar `page`

### Tipos de Consultas Soportadas
- SELECT
- INSERT
- UPDATE
- DELETE
- Stored Procedures

### Limitaciones
- Máximo 1000 filas por página
- Timeout: 5 minutos
- No se permiten múltiples consultas
- Operaciones restringidas:
  - DROP DATABASE
  - CREATE DATABASE
  - BACKUP/RESTORE
  - xp_cmdshell

## Gestión de Usuarios

### Crear Usuario
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nuevo@ejemplo.com",
    "password": "contraseña",
    "role": "user"
  }'
```

### Modificar Usuario
```bash
curl -X PUT http://localhost:8080/api/users/<id> \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "actualizado@ejemplo.com"
  }'
```

## Roles y Permisos

### Roles Disponibles
1. admin
   - Acceso total
   - Gestión de usuarios
   - Todas las operaciones
2. manager
   - Crear conexiones
   - Ejecutar consultas
   - Ver historiales
3. user
   - Ejecutar consultas
   - Ver propio historial

### Asignar Rol
```bash
curl -X POST http://localhost:8080/api/users/<id>/roles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "manager"
  }'
```

## Troubleshooting

### Errores Comunes

1. "Invalid token"
   - Token expirado
   - Token mal formado
   - Solución: Re-autenticarse

2. "Connection failed"
   - Credenciales incorrectas
   - Servidor no disponible
   - Puerto bloqueado
   - Solución: Verificar configuración

3. "Query timeout"
   - Consulta muy larga
   - Servidor sobrecargado
   - Solución: Optimizar consulta

4. "Permission denied"
   - Rol inadecuado
   - Permiso faltante
   - Solución: Contactar administrador

### Contacto Soporte
- Email: support@microsql-ago.com
- Documentación: docs.microsql-ago.com
- Github Issues: github.com/yken-neky/MicroSQL-AGo/issues