# Documentación API - Endpoints Admin

Este documento describe todos los endpoints disponibles para el rol `admin`. Todas las rutas listadas requieren autenticación JWT válida y que el usuario tenga rol `admin` (o el permiso/privilegio que se documenta).

Base: `/api/admin` (ej. `GET /api/admin/users`)

1) Usuarios (admin)
- GET /api/admin/users
  - Descripción: Lista todos los usuarios del sistema junto a los roles asociados.
  - Permiso requerido: `role=admin` (o permiso equivalente).
  - Respuesta 200:
    ```json
    {
      "users": [
        {
          "id": 1,
          "username": "alice",
          "email": "alice@example.com",
          "is_active": true,
          "roles": ["admin","manager"]
        }
      ]
    }
    ```

2) Sessions / tokens
- GET /api/admin/sessions
  - Descripción: lista sesiones/tokens activas (solo visibilidad administrativa).

3) Métricas (ya implementadas)
- GET /api/admin/metrics/users
  - Conteos y distribución por estado/roles
- GET /api/admin/metrics/connections
  - Información sobre conexiones activas y logs
- GET /api/admin/metrics/audits
  - Conteos, duración promedio, estados
- GET /api/admin/metrics/roles
  - Conteo de roles, permisos, distribución role->permission
- GET /api/admin/metrics/system
  - Conteo de filas en tablas críticas y resumen

4) Gestión de Roles y Permisos (CRUD + asignaciones)
- Roles:
  - POST /api/admin/roles — crear rol (payload: { "name": "...", "description": "..." })
  - GET /api/admin/roles — listar roles
  - PUT /api/admin/roles/:id — actualizar
  - DELETE /api/admin/roles/:id — eliminar
  - POST /api/admin/users/:id/roles — asignar rol a usuario (payload: { "role_id": <id> })
  - DELETE /api/admin/users/:id/roles/:role_id — revocar rol

- Permisos:
  - POST /api/admin/permissions — crear permiso (payload detallado en ROLES_PERMISSION_DOC.md)
  - GET /api/admin/permissions — listar permisos
  - PUT /api/admin/permissions/:id — actualizar permiso
  - DELETE /api/admin/permissions/:id — eliminar permiso
  - POST /api/admin/roles/:id/permissions — asignar permiso a rol (payload: { "permission_id": <id> })
  - DELETE /api/admin/roles/:id/permissions/:permission_id — revocar permiso del rol

5) Auditoría de RBAC
- GET /api/admin/audit/rbac?actor_id=<id>&target_type=<type>
  - Lista eventos administrativos (creaciones/actualizaciones de roles/permiso/asignaciones)
  - El sistema guarda entradas append-only en admin_action_logs para trazabilidad.

Seguridad:
- Recomendado: operaciones CRUD y asignaciones deben estar protegidas con rol `admin` o permisos granulares explicitly documentados.
- Auditoría: todo cambio en roles/permissions debe dejar registro en `admin_action_logs`.
