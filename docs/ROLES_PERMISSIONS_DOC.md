# Roles y Permisos — Guía detallada (español)

Este documento explica cómo funcionan roles y permisos, cómo crear roles/permissions nuevos, asignarlos a usuarios y cómo proteger rutas con ellos.

1) Modelo conceptual
- Role (rol): agrupación de permisos.
  - Campos principales: id, name, description.
- Permission (permiso): unidad de privilegio que controla una acción sobre un recurso.
  - Campos principales: id, name, resource, action, description.
  - Ejemplo: `{ "name": "audits:view", "resource": "audits", "action": "view" }`
- UserRole: relación N:M entre users y roles.
- RolePermission: relación N:M entre roles y permissions.

2) Significado de `resource`, `action` y `name`
- resource: dominio que protege, p.ej. `audits`, `users`, `connections`, `metrics`.
- action: verbo/operación: `view`, `create`, `update`, `delete`, `owner:update`.
- name: identificador único y recomendado con formato `resource:action` (ej. `audits:view`).
- Recomendación: usa siempre resource + action para evitar ambigüedades.

3) Ejemplo real — crear `manager:view_audit`
- JSON para crear permiso:
  ```json
  {
    "name": "manager:view_audit",
    "resource": "audits",
    "action": "view",
    "description": "Permite al rol manager ver auditorías"
  }
  ```
- Asignar a rol `manager` (supongamos role_id=7):
  POST /api/admin/roles/7/permissions
  payload: `{ "permission_id": <id_del_permiso> }`

- Asignar rol a usuario:
  POST /api/admin/users/42/roles
  payload: `{ "role_id": 7 }`

4) Protegiendo rutas con permisos y roles
- Recomendado:
  - Uso de middleware por permiso: `authzMW.RequirePermission("audits:view")`
  - Uso de middleware por rol (más simple): `authMW.RequireRole("admin")`
  - Para permisos tipo `owner` (p.ej. `audits:owner:view`): middleware solo valida existencia del permiso; la comprobación de propiedad (que el usuario sea dueño del recurso) debe implementarla el handler.

Handler sample (ownership):
```go
// Pseudocódigo - handler que protege lectura por owner o permiso global
func (h *AuditHandler) GetAudit(c *gin.Context) {
    userID := getUserIDFromContext(c)
    auditID := parseParamID(c, "id")
    audit := h.auditRepo.GetByID(auditID)
    if audit.UserID != userID {
        // si no es propietario, permitir solo si tiene permiso global
        if !h.hasPermission(userID, "audits:view") {
           c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":"forbidden"})
           return
        }
    }
    // devolver recurso
}
```

5) Best practices / recomendaciones
- Denegar por defecto.
- Mantener permisos lo más finos posible: evita permisos globales estilo `admin:*` salvo para super-admins.
- Documentar permisos y roles en un YAML/JSON seed que se aplique en migraciones (bootstrap).
- Registrar auditoría de todos los cambios a roles/permissions (admin_action_logs).
- Tests automatizados: casos owner-vs-global-permission, asignación y revocación.

6) Ejemplos de operaciones administrativas (curl)
- Crear permiso:
  curl -X POST -H "Authorization: Bearer <admin>" -d '{"name":"audits:view","resource":"audits","action":"view","description":"Ver auditorías"}' /api/admin/permissions
- Asignar a rol:
  curl -X POST -H "Authorization: Bearer <admin>" -d '{"permission_id": 10}' /api/admin/roles/3/permissions
- Asignar rol a usuario:
  curl -X POST -H "Authorization: Bearer <admin>" -d '{"role_id": 3}' /api/admin/users/42/roles

7) Notas sobre el campo `name` de permisos
- `name` puede ser usado por el middleware como llave rápida y por UIs. Mantener formato `resource:action` permite convenciones claras.
- Evitar duplicados en `name`. El sistema asume unicidad.

---

Si quieres, implemento ahora:
- validación extra en los handlers que use `owner` de forma estándar (middleware reusable), y/o
- añadir ejemplos de seed en formato migration que creen un conjunto base de roles/permissions.

¿Quieres que añada el middleware reutilizable que valida ownership y un ejemplo en un endpoint concreto?// filepath: docs/ROLES_PERMISSION_DOC.md
# Roles y Permisos — Guía detallada (español)

Este documento explica cómo funcionan roles y permisos, cómo crear roles/permissions nuevos, asignarlos a usuarios y cómo proteger rutas con ellos.

1) Modelo conceptual
- Role (rol): agrupación de permisos.
  - Campos principales: id, name, description.
- Permission (permiso): unidad de privilegio que controla una acción sobre un recurso.
  - Campos principales: id, name, resource, action, description.
  - Ejemplo: `{ "name": "audits:view", "resource": "audits", "action": "view" }`
- UserRole: relación N:M entre users y roles.
- RolePermission: relación N:M entre roles y permissions.

2) Significado de `resource`, `action` y `name`
- resource: dominio que protege, p.ej. `audits`, `users`, `connections`, `metrics`.
- action: verbo/operación: `view`, `create`, `update`, `delete`, `owner:update`.
- name: identificador único y recomendado con formato `resource:action` (ej. `audits:view`).
- Recomendación: usa siempre resource + action para evitar ambigüedades.

3) Ejemplo real — crear `manager:view_audit`
- JSON para crear permiso:
  ```json
  {
    "name": "manager:view_audit",
    "resource": "audits",
    "action": "view",
    "description": "Permite al rol manager ver auditorías"
  }
  ```
- Asignar a rol `manager` (supongamos role_id=7):
  POST /api/admin/roles/7/permissions
  payload: `{ "permission_id": <id_del_permiso> }`

- Asignar rol a usuario:
  POST /api/admin/users/42/roles
  payload: `{ "role_id": 7 }`

4) Protegiendo rutas con permisos y roles
- Recomendado:
  - Uso de middleware por permiso: `authzMW.RequirePermission("audits:view")`
  - Uso de middleware por rol (más simple): `authMW.RequireRole("admin")`
  - Para permisos tipo `owner` (p.ej. `audits:owner:view`): middleware solo valida existencia del permiso; la comprobación de propiedad (que el usuario sea dueño del recurso) debe implementarla el handler.

Handler sample (ownership):
```go
// Pseudocódigo - handler que protege lectura por owner o permiso global
func (h *AuditHandler) GetAudit(c *gin.Context) {
    userID := getUserIDFromContext(c)
    auditID := parseParamID(c, "id")
    audit := h.auditRepo.GetByID(auditID)
    if audit.UserID != userID {
        // si no es propietario, permitir solo si tiene permiso global
        if !h.hasPermission(userID, "audits:view") {
           c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":"forbidden"})
           return
        }
    }
    // devolver recurso
}
```

5) Best practices / recomendaciones
- Denegar por defecto.
- Mantener permisos lo más finos posible: evita permisos globales estilo `admin:*` salvo para super-admins.
- Documentar permisos y roles en un YAML/JSON seed que se aplique en migraciones (bootstrap).
- Registrar auditoría de todos los cambios a roles/permissions (admin_action_logs).
- Tests automatizados: casos owner-vs-global-permission, asignación y revocación.

6) Ejemplos de operaciones administrativas (curl)
- Crear permiso:
  curl -X POST -H "Authorization: Bearer <admin>" -d '{"name":"audits:view","resource":"audits","action":"view","description":"Ver auditorías"}' /api/admin/permissions
- Asignar a rol:
  curl -X POST -H "Authorization: Bearer <admin>" -d '{"permission_id": 10}' /api/admin/roles/3/permissions
- Asignar rol a usuario:
  curl -X POST -H "Authorization: Bearer <admin>" -d '{"role_id": 3}' /api/admin/users/42/roles

7) Notas sobre el campo `name` de permisos
- `name` puede ser usado por el middleware como llave rápida y por UIs. Mantener formato `resource:action` permite convenciones claras.
- Evitar duplicados en `name`. El sistema asume unicidad.

---

Si quieres, implemento ahora:
- validación extra en los handlers que use `owner` de forma estándar (middleware reusable), y/o
- añadir ejemplos de seed en formato migration que creen un conjunto base de roles/permissions.

¿Quieres que añada el middleware reutilizable que valida ownership y un ejemplo en un endpoint concreto?