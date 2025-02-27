# Bienvenido a MicroSQL AGo (v0.1)

> [!NOTE] 
> 
> Un microservicio para las auditorías de seguridad automatizadas en Microsoft SQL Server.

---

#### Estructura del proyecto: 

```
MicroSQL Ago/
│
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│       ├── models/
│       │   └── postgresModels.go
|	    ├──	connections/
│       │   └── postgresConnections.go
│       └── controllers/
│           └── postgresControllers.go
├── remote/
│       ├── services/
│       │   └── remoteServices.go
|	    ├──	connections/
│       │   └── remoteConnections.go
│       └── utils/
│           └── utils.go
├── config/
│   └── deploy/ #todos los .yml para el despliegue en kubernetes
│       └── deployment.yml
├── routes/
│   └── routes.go
├── vendor/
│   └── # dependencias creadas por "go mod vendor"    
├── .env
├── .gitignore
├── go.mod
│   └── go.sum
└── README.md
```