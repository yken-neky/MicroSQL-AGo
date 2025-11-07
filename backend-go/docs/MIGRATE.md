# Migración de datos desde SQLite a MySQL

Este documento explica cómo usar el migrador incluido en el proyecto para copiar los datos existentes de la base de datos SQLite local hacia una base de datos MySQL.

Hay dos formas principales de ejecutar el migrador:

1. Localmente (en tu máquina de desarrollo) con `go run`.
2. Dentro de Docker (usando `docker compose`), montando el archivo `db.sqlite3` si existe en tu máquina.

---

## Requisitos

- Go (para ejecuciones locales) o Docker (para ejecutar con contenedores).
- El archivo SQLite de origen (por defecto `./db.sqlite3`) debe existir y contener las tablas que produce el proyecto.
- Un servidor MySQL accesible y las variables de entorno configuradas: `MYSQL_HOST`, `MYSQL_USER`, `MYSQL_PASSWORD`, `MYSQL_DATABASE`. Alternativamente, puedes pasar el DSN de MySQL vía `--dst-dsn`.

## Uso local (ejemplo)

1. Asegúrate de tener el módulo y dependencias:

```powershell
cd D:\Proyectos\MicroSQL AGo\backend-go
go mod tidy
```

2. Ejecuta el migrador indicando la ruta al archivo SQLite si no es `./db.sqlite3`.

Ejemplo (usando variables de entorno para MySQL):

```powershell
setx MYSQL_HOST "127.0.0.1"
setx MYSQL_PORT "3306"
setx MYSQL_USER "micro_user"
setx MYSQL_PASSWORD "ChangeMeStrongPassword!"
setx MYSQL_DATABASE "microsql_ago"

go run ./cmd/migrate --src ./db.sqlite3
```

Ejemplo (proporcionando DSN completo en la línea de comandos):

```powershell
go run ./cmd/migrate --src ./db.sqlite3 --dst-dsn "user:pass@tcp(127.0.0.1:3306)/microsql_ago?charset=utf8mb4&parseTime=True&loc=Local"
```

## Uso en Docker

Por defecto `docker-compose.yml` incluye un servicio `migrate` que ejecuta `/app/migrate`, pero para que pueda leer tu archivo SQLite local deberás montarlo como volumen en el servicio `migrate`.

Ejemplo de `docker-compose` (fragmento) para montar el archivo local `db.sqlite3`:

```yaml
  migrate:
    build: .
    depends_on:
      - db
    env_file:
      - .env
    volumes:
      - ./db.sqlite3:/app/db.sqlite3:ro
    command: ["/app/migrate"]
    restart: "no"
```

Con ese montaje, si `DB_PATH` en tu `.env` apunta a `/app/db.sqlite3` (o usas `--src /app/db.sqlite3`), el migrador podrá leer tu archivo y copiar los datos a MySQL.

También puedes ejecutar el migrador a mano usando `docker compose run` y montando el archivo:

```powershell
docker compose build migrate
docker compose run --rm -v %cd%/db.sqlite3:/app/db.sqlite3 migrate /app/migrate --src /app/db.sqlite3
```

## Notas y precauciones

- El migrador realiza inserciones usando `OnConflict DoNothing` para evitar insertar duplicados si las filas ya existen en la BD destino.
- Revisa los tipos y constraints después de migrar; GORM `AutoMigrate` intenta crear tablas compatibles pero podría necesitar ajustes manuales para índices o columnas complejas.
- No borra ni modifica datos en la BD destino, sólo intenta insertar los registros encontrados en la fuente.
- En entornos de producción, realiza un backup antes de la migración.

## Ejemplo rápido

1. Montar MySQL local o via `docker compose`.
2. Ejecutar migración local:

```powershell
go run ./cmd/migrate --src ./db.sqlite3
```

o con DSN explícito:

```powershell
go run ./cmd/migrate --src ./db.sqlite3 --dst-dsn "user:pass@tcp(host:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
```

---

Si quieres, puedo añadir soporte para excluir tablas, migrar en batches, o un modo de `dry-run` que solo muestre cuántas filas se migrarían sin ejecutar inserts. ¿Quieres que añada alguna de estas funciones? 
