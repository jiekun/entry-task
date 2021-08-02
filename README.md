# Entry Task
## Introduction
`entry-task` is a web service based on **Gin** framework. It provides basic user management feature.

## Project Layout
```
├── LICENSE
├── README.md
├── cmd
│    ├── grpc-server                        // 
│    └── http-server                        // 
│        └── main.go                        // HTTP Server entry point
├── configs                                 // 
│    └── config.yaml                        // 
├── global                                  // 
│    ├── cache.go                           // 
│    ├── db.go                              // 
│    └── setting.go                         // 
├── go.mod                                  // 
├── go.sum                                  // 
├── internal                                // Service Codes
│    ├── constant                           // 
│    ├── dao                                // Data Access Object Layer
│    ├── error                              // 
│    ├── models                             // 
│    │    ├── model.go                      // 
│    │    └── user.go                       // 
│    ├── routers                            // 
│    │    ├── api                           // Controller Layer for input params handling
│    │    │    ├── ping.go                  // 
│    │    │    ├── upload.go                // 
│    │    │    └── user.go                  // 
│    │    └── router.go                     // 
│    └── service                            // Service Logic Layer
│        ├── service.go                     // 
│        ├── upload.go                      // 
│        └── user.go                        // 
├── log                                     // 
├── pkg                                     // Public utils
│    ├── hashing                            // 
│    ├── logger                             // 
│    ├── middleware                         // 
│    ├── resp                               // 
│    ├── setting                            // 
│    └── upload                             // 
├── scripts                                 // Scripts for initialization 
└── upload                                  // User upload contents
```

## Setup
```
go mod tidy

go run cmd/http-server/main.go
```