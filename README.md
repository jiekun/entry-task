# Entry Task
## Introduction
`entry-task` is a web service based on **Gin** framework. It provides basic user management feature.

## Project Layout
```
├── LICENSE
├── README.md
├── cmd
│   ├── grpc-server
│   └── http-server
│       └── main.go
├── configs
├── go.mod
├── go.sum
└── internal
    └── routers
        ├── api
        │   └── ping.go
        └── router.go
```

## Setup
```
go mod tidy

go run cmd/http-server/main.go
```