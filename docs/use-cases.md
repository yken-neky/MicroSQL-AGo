# Casos de Uso

## 1. Gestión de Conexiones

### 1.1 Establecer Conexión a SQL Server
**Actor Principal:** Usuario autenticado
**Precondiciones:**
- Usuario autenticado con token válido
- Permisos de conexión asignados

**Flujo Principal:**
1. Usuario envía credenciales de SQL Server
2. Sistema valida credenciales
3. Sistema establece conexión
4. Sistema registra conexión activa
5. Sistema retorna detalles de conexión

**Flujos Alternativos:**
- Credenciales inválidas: Sistema notifica error
- Límite de conexiones alcanzado: Sistema notifica error
- Servidor no disponible: Sistema reintenta y notifica error

### 1.2 Cerrar Conexión
**Actor Principal:** Usuario autenticado
**Precondiciones:**
- Conexión activa existente

**Flujo Principal:**
1. Usuario solicita cerrar conexión
2. Sistema cierra conexión física
3. Sistema actualiza registro
4. Sistema libera recursos

## 2. Ejecución de Consultas

### 2.1 Ejecutar Consulta SELECT
**Actor Principal:** Usuario autenticado
**Precondiciones:**
- Conexión activa
- Permisos de ejecución

**Flujo Principal:**
1. Usuario envía consulta SQL
2. Sistema valida sintaxis
3. Sistema ejecuta consulta
4. Sistema procesa resultados
5. Sistema retorna datos paginados

**Flujos Alternativos:**
- Sintaxis inválida: Sistema notifica error
- Timeout: Sistema cancela ejecución
- Error de ejecución: Sistema registra error

### 2.2 Ejecutar Consulta de Modificación
**Actor Principal:** Usuario autenticado con permisos elevados
**Precondiciones:**
- Conexión activa
- Permisos de modificación

**Flujo Principal:**
1. Usuario envía consulta INSERT/UPDATE/DELETE
2. Sistema valida operación
3. Sistema ejecuta consulta
4. Sistema retorna filas afectadas

## 3. Gestión de Usuarios

### 3.1 Crear Usuario
**Actor Principal:** Administrador
**Precondiciones:**
- Usuario con rol admin

**Flujo Principal:**
1. Admin envía datos de nuevo usuario
2. Sistema valida datos
3. Sistema crea usuario
4. Sistema asigna rol inicial

### 3.2 Asignar Roles
**Actor Principal:** Administrador
**Precondiciones:**
- Usuario existente
- Rol válido

**Flujo Principal:**
1. Admin selecciona usuario y rol
2. Sistema verifica permisos
3. Sistema asigna rol
4. Sistema actualiza permisos

## 4. Monitoreo y Auditoría

### 4.1 Consultar Historial
**Actor Principal:** Usuario autenticado
**Precondiciones:**
- Usuario con permisos de lectura

**Flujo Principal:**
1. Usuario solicita historial
2. Sistema filtra por usuario/fecha
3. Sistema retorna registros paginados

### 4.2 Monitorear Rendimiento
**Actor Principal:** Administrador
**Precondiciones:**
- Acceso a métricas

**Flujo Principal:**
1. Admin accede a dashboard
2. Sistema recopila métricas
3. Sistema muestra estadísticas

## 5. Gestión de Caché

### 5.1 Gestionar Caché de Consultas
**Actor Principal:** Sistema
**Precondiciones:**
- Redis disponible

**Flujo Principal:**
1. Sistema recibe consulta
2. Sistema verifica caché
3. Sistema retorna resultado o ejecuta
4. Sistema actualiza caché

## 6. Seguridad

### 6.1 Validar Operaciones Sensibles
**Actor Principal:** Sistema
**Precondiciones:**
- Consulta recibida

**Flujo Principal:**
1. Sistema analiza consulta
2. Sistema verifica palabras clave
3. Sistema valida permisos
4. Sistema permite/deniega

## 7. Alta Disponibilidad

### 7.1 Gestionar Reconexión
**Actor Principal:** Sistema
**Precondiciones:**
- Pérdida de conexión detectada

**Flujo Principal:**
1. Sistema detecta desconexión
2. Sistema intenta reconexión
3. Sistema notifica estado
4. Sistema restaura operación