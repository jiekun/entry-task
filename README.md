# Entry Task
## Introduction
`entry-task` is a web service based on **Gin** framework. It provides basic user management feature.

## Project Layout
```
├── LICENSE
├── README.md
├── cmd
│   ├── grpc-server
│   └── http-server
│       └── main.go           // Entry point
├── configs
│   └── config.yaml           // Config file
├── global
├── go.mod
├── go.sum
├── internal
│   ├── constant
│   ├── dao
│   ├── error                 // Error handling and definition
│   ├── models
│   ├── routers
│   │   ├── api               // URL router
│   │   └── router.go
│   └── service
├── log
├── pkg
│   ├── logger
│   ├── resp
│   └── setting
└── scripts                   // All init scripts and sql files
```

## Setup
```
go mod tidy

go run cmd/http-server/main.go
```