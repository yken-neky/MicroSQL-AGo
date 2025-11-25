# Logging & Observability - Guidelines

Fecha: 2025-11-25

Este documento describe el enfoque y el uso del sistema de logging y pequeñas observability helpers implementados en el proyecto.

## Objetivos

- Proveer un `zap` logger central y consistente para todas las capas.
- Añadir un middleware que adjunte un logger por request con `request_id` para correlación.
- Registrar queries de GORM con logger estructurado (slow query + errores).
- Redactar valores sensibles (tokens/passwords) antes de exponerlos en los logs.

## Formato y contractos

- Todos los logs están estructurados (JSON compatible a través de `zap`).
- Campos mínimos por request logger:
  - request_id (UUID)
  - method
  - path
  - remote_ip
  - status
  - duration_ms

## Archivos clave

- `internal/adapters/primary/http/middleware/logging_middleware.go` — middleware que crea el request logger y lo atacha al `gin.Context` usando la key `logger`.
- `internal/pkg/logging/redact.go` — funciones utilitarias para redacción de Authorization tokens y valores sensibles.
- `internal/infrastructure/persistence/gorm_logger.go` — adaptador zap -> gorm/logger.Interface para logging de SQL (incluye slow-query detection y errores).
- `cmd/server/main.go` — el servidor ahora instala el logger de GORM en la sesión del DB para que SQL se registre con zap.

## Cómo usar desde handlers

Dentro de un handler Gin puedes obtener el request-scoped logger así:

```go
lg, ok := c.Get("logger")
if ok {
    zapLogger := lg.(*zap.Logger)
    zapLogger.Info("handling request", zap.String("someKey", "value"))
}
```

Si no existe `logger` usa el `pkg/utils.NewLogger` que se crea en `cmd/server`.

## Redacción y seguridad

- Authorization headers y tokens se redaccionan usando `RedactAuthHeader`.
- Redactar passwords antes de loguear valores de entrada (no registrar contraseñas en claro).

## GORM logging

- SQL logs se emiten con `zap` mediante `internal/infrastructure/persistence/gorm_logger.go`.
- Slow threshold por defecto: 200ms (configurable en creación del logger en `main.go`).

## Test y verificación

- Hay tests unitarios básicos para middleware y redactor (`*_test.go`) en el árbol del proyecto.

## Recomendaciones futuras

- Añadir correlación con tracing (OpenTelemetry) para propagar `trace_id` y exportar traces.
- En production, configurar un sink (Loki/ELK/CloudWatch) y retention.
- Añadir sampling para no inundar logs en queries frecuentes de debug.
