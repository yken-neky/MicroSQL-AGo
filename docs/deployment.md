# Documentación de Despliegue

## Requisitos del Sistema

### Producción
- Linux/Windows Server 2019+
- 4GB RAM mínimo
- 2 CPUs
- 20GB espacio en disco
- Conexión a Internet

### Software Requerido
- Go 1.20+
- SQL Server 2019+
- Redis 6+
- Docker (opcional)
- Nginx (recomendado para producción)

## Instalación

### 1. Preparación del Entorno

```bash
# Crear directorios
mkdir -p /opt/microsql-ago/{configs,logs,data}
cd /opt/microsql-ago

# Clonar repositorio
git clone https://github.com/yken-neky/MicroSQL-AGo.git .

# Copiar configuración
cp configs/env.example configs/.env
```

### 2. Configuración de Base de Datos

```bash
# SQL Server
docker run -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=YourStrong!Passw0rd' \
   -p 1433:1433 --name sqlserver \
   -d mcr.microsoft.com/mssql/server:2019-latest

# Redis
docker run --name redis -p 6379:6379 -d redis:6
```

### 3. Variables de Entorno

Editar `.env`:
```env
APP_PORT=8080
APP_ENV=production
LOG_LEVEL=info

SQLSERVER_HOST=localhost
SQLSERVER_PORT=1433
SQLSERVER_USER=sa
SQLSERVER_PASSWORD=YourStrong!Passw0rd

REDIS_HOST=localhost
REDIS_PORT=6379
```

### 4. Compilación

```bash
# Instalar dependencias
go mod download

# Compilar
make build
```

### 5. Configuración del Servicio

Crear `/etc/systemd/system/microsql-ago.service`:
```ini
[Unit]
Description=MicroSQL AGo Service
After=network.target

[Service]
Type=simple
User=microsql
WorkingDirectory=/opt/microsql-ago
ExecStart=/opt/microsql-ago/build/microsql-ago
Restart=always
Environment=APP_ENV=production

[Install]
WantedBy=multi-user.target
```

### 6. Despliegue con Docker

```bash
# Construir imagen
docker build -t microsql-ago:latest .

# Ejecutar con docker-compose
docker-compose up -d
```

## Monitorización

### Logs
- Aplicación: `/opt/microsql-ago/logs/app.log`
- Access Log: `/opt/microsql-ago/logs/access.log`
- Error Log: `/opt/microsql-ago/logs/error.log`

### Métricas
- Prometheus endpoint: `:8080/metrics`
- Grafana Dashboard incluido en `./configs/grafana/`

### Healthcheck
- Endpoint: `/health`
- Prometheus: `/metrics`

## Backup y Recuperación

### Base de Datos
```bash
# Backup
./scripts/backup.sh

# Restore
./scripts/restore.sh <backup_file>
```

### Configuración
```bash
# Backup configs
tar -czf config-backup.tar.gz configs/

# Restore configs
tar -xzf config-backup.tar.gz
```

## Escalamiento

### Horizontal
1. Configurar Redis para sesiones distribuidas
2. Configurar balanceador de carga
3. Ajustar pool de conexiones

### Vertical
1. Aumentar recursos del servidor
2. Ajustar configuración de memoria
3. Optimizar índices de BD

## Troubleshooting

### Problemas Comunes

1. Error de Conexión a SQL Server
```
Error: dial tcp connection refused
```
Solución: Verificar firewall y credenciales

2. Redis no disponible
```
Error: redis connection failed
```
Solución: Verificar servicio Redis y conectividad

3. Memoria Insuficiente
```
Error: out of memory
```
Solución: Ajustar límites de memoria en `.env`

### Comandos Útiles

```bash
# Verificar estado
systemctl status microsql-ago

# Ver logs
journalctl -u microsql-ago

# Reiniciar servicio
systemctl restart microsql-ago
```

## Seguridad

### Firewall
```bash
# Abrir puertos necesarios
ufw allow 8080/tcp  # API
ufw allow 1433/tcp  # SQL Server
ufw allow 6379/tcp  # Redis
```

### SSL/TLS
1. Generar certificados
2. Configurar Nginx
3. Actualizar configuración

### Hardening
1. Deshabilitar acceso root
2. Configurar fail2ban
3. Implementar WAF