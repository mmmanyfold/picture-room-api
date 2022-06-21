# `picture-room-api`

> API Server for Picture Room V2

Project layout based on [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

### Project Layout

```
{root}

├── internal # Private application and library code
├── cmd      # Main driver for Golang app (our entry point)
├── pkg      # Public library code to be used by external applications
    └── api  # API controller business logic
    └── db   # Database connection and models (lightweight schemas provided by gorm)
...
```

### Requirements

- [golang](https://golang.org/)
- [docker](https://docs.docker.com/get-docker/)
- [docker-compose](https://docs.docker.com/compose/install/)
- sqlc - `brew install kyleconroy/sqlc/sqlc`
- (optional) [goose for migrations](https://github.com/pressly/goose)
- (optional) [google cloud sdk for deployment](https://cloud.google.com/sdk/docs/downloads-versioned-archives)

### Setup

##### create alias for docker mac
- `sudo ifconfig lo0 alias 10.254.254.254` 

##### or linux
- `sudo apt-get install net-tools -y`
- `sudo ifconfig lo:0 10.254.254.254.254`

##### setup configs
- `cp sample.config.yaml config/dev.yaml`
- update `config/dev.yaml` with development values

### Development

- touch `/config/prod.yaml` and update it with deployment values
- `make up`

### Tests

- `make test` # runs unit tests
- `make test-ci` # runs integration tests

### Deploy

- `make deploy` # deploys app to google cloud engine

### Migrations

With goose installed (`go get -u github.com/pressly/goose/cmd/goose`)

##### create new migration

`goose -dir internal/db/migrations create <migration_name> sql`

##### run migrations

- `make migrate-up`

- `make migrate-down`
