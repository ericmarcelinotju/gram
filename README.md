# GRAM

GRAM is Golang boilerplate designed to create neat REST API with robust way to design multiple repository from one interface.

![Project structure](./diagram.png?raw=true "Project Structure")

# Features

- Migration
- Seeding
- Multiple Database Driver :
  - MySQL
  - SQLite
  - SQLServer
  - PostgreSQL
- Scheduler
- Websocket (need message queue)

# Development

```
go run main.go
```

# Commands

> Create super user
```
go run main.go -u [Super User Name]
```

> Migrate
```
go run main.go -m
```

> Seeding
```
go run main.go -s
```

# Swagger

All routes registered at Routes/Controller directory can be documented using gin-swagger comment

> Update `swagger.json` and `swagger.yaml`
```
swag init
```

> Open swagger documentation page
```
{{baseUrl}}/swagger/index.
```