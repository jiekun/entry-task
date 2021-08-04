# Entry Task
## Introduction
`entry-task` is a web service based on **Gin** framework. It provides basic user management feature.

## Document
You could find following guidance on [Entry Task System Design Document](https://docs.google.com/document/d/1sd5S8xdJRYcZrYAOM1cREnuQslZnkj8kIIQccHNmlq4/edit#) and `README.md`:
- Feature Description
- API Doc
- Deploy Doc
- Benchmark Result

## Setup
Before getting start, you need **MySQL** and **Redis** installed on your machine. And please refer `go.mod` for the Go packages used in this project.`init.sql` in script folder can help to initialize database & table we need.

Following commands will help you set up the project and run it locally:
```bash
# generate a config file
cp ./config/config.yaml.defaul ./config/config.yaml

# modify the host/db/etc settings in config file
vim ./config/config.yaml

# packages dependancy
go mod tidy

# run a RPC server
go run ./cmd/grpc-server/main.go

# run a HTTP server
go run ./cmd/http-server/main.go 
```

Lastly try to test if it works with:
```bash
curl --location --request GET 'http://{your_http_host}/api/ping' # {"message":"pong"}
```

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