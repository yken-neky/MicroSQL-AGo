# TODO List para completar el servicio al 100%

Este documento lista las tareas pendientes organizadas por área funcional y prioridad para llevar el backend a producción. Incluye breves criterios de aceptación y estimaciones orientativas.

Formato de cada ítem:
- Prioridad: Alta / Media / Baja
- Est. esfuerzo: horas (aprox.)
- Owner: TBD
- Criterio de aceptación: cómo verificar que está resuelto

---

## 1. Infraestructura y despliegue

- Configurar y probar despliegue con Docker Compose (MySQL + app + migrate)
  - Prioridad: Alta
  - Est. esfuerzo: 2h
  - Owner: TBD
  - Criterio de aceptación: `docker compose up --build` levanta DB y app; migrator puede ejecutarse y crear tablas.

- Añadir healthchecks y readiness probes para contenedores
  - Prioridad: Alta
  - Est. esfuerzo: 3h
  - Criterio de aceptación: `docker compose ps` muestra contenedores healthy y app no recibe tráfico hasta ready.

- Configurar CI/CD: build, lint, tests, image push
  - Prioridad: Alta
  - Est. esfuerzo: 8h
  - Criterio de aceptación: Pipeline verde con build y tests, imagen publicada en registry.

---

## 2. Persistencia y migraciones

- Soporte completo MySQL en producción (DSN, pool, TLS)
  - Prioridad: Alta
  - Est. esfuerzo: 4h
  - Criterio de aceptación: App se conecta a MySQL con pool configurado y TLS opcional.

- Revisar y ajustar migraciones GORM (tipos, índices, constraints)
  - Prioridad: Alta
  - Est. esfuerzo: 4h
  - Criterio de aceptación: `cmd/migrate` crea tablas correctamente en MySQL; `SHOW CREATE TABLE` OK.

- Script/flow seguro para migrar datos desde SQLite a MySQL (idempotente, batches)
  - Prioridad: Alta
  - Est. esfuerzo: 6h
  - Criterio de aceptación: migración ejecutable sin duplicados y con opción `--dry-run`.

- Configurar backups y restauración (mysqldump o snapshot)
  - Prioridad: Alta
  - Est. esfuerzo: 6h
  - Criterio de aceptación: Procedimiento documentado y script de backup funcional.

---

## 3. Autenticación y autorización

- Implementar login real (verificar credenciales, JWT signing)
  - Prioridad: Alta
  - Est. esfuerzo: 6h
  - Criterio de aceptación: `POST /api/users/login` devuelve JWT válido; expiración y claims ok.

- Implementar middleware de autenticación (JWT) y autorización por roles
  - Prioridad: Alta
  - Est. esfuerzo: 6h
  - Criterio de aceptación: rutas protegidas retornan 401/403 correctamente; tests unitarios.

- Gestión de usuarios: crear/editar/desactivar (with password hashing)
  - Prioridad: Alta
  - Est. esfuerzo: 8h
  - Criterio de aceptación: endpoints `/api/users` crean y actualizan usuarios con contraseñas hasheadas.

---

## 4. API: Endpoints de negocio (implementación real)

- Conexiones SQL (connect/disconnect/list/detail)
  - Prioridad: Alta
  - Est. esfuerzo: 12h
  - Criterio de aceptación: `POST /api/connections` crea conexión (persistida cifrada), `GET /api/connections` lista, `DELETE /api/connections/{id}` cierra.

- Ejecución de consultas (async/streaming/pagination)
  - Prioridad: Alta
  - Est. esfuerzo: 20h
  - Criterio de aceptación: `POST /api/queries/execute` ejecuta SQL con paginación, respeta límites de concurrencia y rate limits.

- Historial de consultas y resultados (`/api/queries/history`, `/api/queries/{id}`)
  - Prioridad: Alta
  - Est. esfuerzo: 8h
  - Criterio de aceptación: historial persistido y accesible con filtros y paginación.

- Controles / Auditorías (list, execute, result)
  - Prioridad: Media
  - Est. esfuerzo: 12h
  - Criterio de aceptación: endpoints para listar controles y ejecutar auditorías parciales/complete.

- Dashboard endpoints (audits amount, connections amount, correct rate)
  - Prioridad: Media
  - Est. esfuerzo: 6h
  - Criterio de aceptación: métricas calculadas y endpoints devuelven JSON válidos.

---

## 5. Seguridad y hardening

- Hashing seguro de contraseñas (bcrypt/scrypt/argon2)
  - Prioridad: Alta
  - Est. esfuerzo: 4h
  - Criterio de aceptación: contraseñas almacenadas hasheadas; tests.

- Protección contra inyección SQL (sanitizar inputs, parametrizar queries)
  - Prioridad: Alta
  - Est. esfuerzo: 8h
  - Criterio de aceptación: pruebas de no-inyección y revisión de QueryExecutor.

- Rate limiting y control de concurrencia por usuario/IP
  - Prioridad: Alta
  - Est. esfuerzo: 6h
  - Criterio de aceptación: límites aplicados y headers X-RateLimit presentes.

- Gestión de secretos (no usar .env en prod; usar vault/secret manager)
  - Prioridad: Alta
  - Est. esfuerzo: 6h
  - Criterio de aceptación: secretos no en repo; doc con instrucciones de integración.

---

## 6. Observabilidad

- Logging estructurado y niveles (zap ya presente)
  - Prioridad: Media
  - Est. esfuerzo: 4h
  - Criterio de aceptación: logs JSON consistentes con campos clave (request id, user id).

- Métricas Prometheus y endpoints `/metrics`
  - Prioridad: Media
  - Est. esfuerzo: 6h
  - Criterio de aceptación: métricas básicas expuestas (request_count, latency, db_connections).

- Tracing distribuido (opcional)
  - Prioridad: Baja
  - Est. esfuerzo: 12h

---

## 7. Tests

- Unit tests para casos de uso y repositorios (coverage objetivo >= 70%)
  - Prioridad: Alta
  - Est. esfuerzo: 24h
  - Criterio de aceptación: tests pasan en CI y coverage report >= threshold.

- Integration tests (usar MySQL en CI, o sqlite con casos reales)
  - Prioridad: Alta
  - Est. esfuerzo: 16h
  - Criterio de aceptación: tests end-to-end que validen flows críticos (login, connect, execute query).

- E2E tests contra docker-compose (smoke tests)
  - Prioridad: Media
  - Est. esfuerzo: 12h

---

## 8. Docs y UX

- Completar `docs/api.md` con payloads reales y ejemplos
  - Prioridad: Alta
  - Est. esfuerzo: 4h

- Generar Swagger/OpenAPI y servir `/api/swagger`
  - Prioridad: Media
  - Est. esfuerzo: 6h

- README con pasos de despliegue y desarrollo (docker, migrate, env)
  - Prioridad: Alta
  - Est. esfuerzo: 3h

---

## 9. Operaciones y producción

- Ajustes de DB pool y performance tuning
  - Prioridad: Alta
  - Est. esfuerzo: 6h

- Backups automáticos y restauración documentada
  - Prioridad: Alta
  - Est. esfuerzo: 6h

- Monitoreo y alertas (uptime, error rate, slow queries)
  - Prioridad: Alta
  - Est. esfuerzo: 8h

---

## 10. Extras y mejoras futuras

- Modo multi-tenant / separación por cliente (si aplica)
- UI administrativa para gestionar usuarios, roles y auditorías
- Soporte para otros RDBMS (Postgres) y migraciones cross-db

---

## Prioridad inmediata (primer sprint sugerido)
1. Autenticación real + middleware JWT
2. Implementar Conexiones (connect/disconnect/list)
3. Implementar migraciones y flujo seguro de datos (migrate)
4. CI básico (build + tests)
5. Completar docs y README

---

Si quieres, puedo:
- Crear issues/tareas separadas en el repo (si usas GitHub/GitLab) por cada ítem.
- Empezar a implementar la prioridad inmediata ahora (elige cuál).
