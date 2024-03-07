# Banking System
- [Banking System](#banking-system)
- [Project structure](#project-structure)
- [Start the server](#start-the-server)
- [Stop the server](#stop-the-server)
- [Add New Restful api](#add-new-restful-api)
    - [Add new route](#add-new-route)
    - [Add new handler for user](#add-new-handler-for-user)
    - [Add new service for user](#add-new-service-for-user)
    - [Add new repo for user using mysql database](#add-new-repo-for-user-using-mysql-database)


# Project structure
```
banking/
├─ app/
│  ├─ api/
│  │  ├─ v1/
│  │  │  ├─ handler/
│  │  │  │  ├─ user/
│  │  │  │  │  ├─ user.go
│  │  │  │  │  ├─ user_test.go
│  │  │  ├─ response.go
│  │  ├─ router.go
│  ├─ repo/
│  │  ├─ mysql/
│  │  │  ├─ user/
│  │  │  │  ├─ query.go
│  │  │  │  ├─ query_test.go
│  │  │  │  ├─ command.go
│  │  │  │  ├─ command_test.go
│  │  │  │  ├─ errorMsg.go
│  │  │  │  ├─ setup_test.go
│  ├─ service/
│  │  ├─ user/
│  │  │  ├─ mock/
│  │  │  │  ├─ user.go
│  │  │  ├─ user.go
├─ build/
│  ├─ docker-compose.yml
├─ cmd/
│  ├─ apiserver.go
│  ├─ root.go
├─ config/
│  ├─ config.example.yaml
├─ database/
│  ├─ mysql/
│  │  ├─ mysql.go
├─ docs/
│  ├─ docs.go
│  ├─ swagger.json
│  ├─ swagger.yaml
├─ global/
│  ├─ global.go
├─ log/
│  ├─ logger.go
├─ model/
│  ├─ user.go
├─ .gitignore
├─ Dockerfile
├─ go.mod
├─ go.sum
├─ LICENSE
├─ main.go
├─ makefile
├─ README.md
```

# Start the server
```bash
make docker_up
```

# Stop the server
```bash
make docker_down
```

# Add New Restful api
### Add new route
1. Add path in [app/api/router.go](app/api/router.go)

### Add new handler for user
1. Add handler in [app/api/v1/handler/user/user.go](app/api/v1/handler/user/user.go)
2. Add handler test in [app/api/v1/handler/user/user_test.go](app/api/v1/handler/user/user_test.go)

### Add new service for user
* The service of this demo project is quite simple, so we don't need to add service test, if more complex logic is needed, we should add service test.
1. Add service in [app/service/user/user.go](app/service/user/user.go)

### Add new repo for user using mysql database
* Query repo is used to get data from database
1. Add query repo in [app/repo/mysql/user/query.go](app/repo/mysql/user/query.go)
2. Add query repo test in [app/repo/mysql/user/query_test.go](app/repo/mysql/user/query_test.go)
* Command repo is used to insert, update, delete data from database
1. Add command repo in [app/repo/mysql/user/command.go](app/repo/mysql/user/command.go)
2. Add command repo test in [app/repo/mysql/user/command_test.go](app/repo/mysql/user/command_test.go)
