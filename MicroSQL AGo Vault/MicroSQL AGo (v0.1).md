> [!ERROR] Un microservicio para las auditorías de seguridad automatizadas en Microsoft SQL Server.

---
# En esta versión (v0.1): 

#### Definida la estructura: 

```
main/
│
├── local/
│   ├── controllers/
│   │   └── dbController.go
│   ├── models/
│   │   └── dbModels.go
|	├──	connections/
│   │   └── dbConnections.go
│   └── routes/
│       └── routes.go
│
└── main.go
```

#### Motor de base de datos para el uso del microservicio en fase de pruebas: **PostgreSQL**

#### Framework para el enrutamiento: **Gin**

---

> [!DONE] 
> #### Actualmente funcionando: 
> - Conexión con la base de datos local
> - Listar objetos de la tabla definida en el modelo
> - Crear objetos de la tabla definida en el modelo

