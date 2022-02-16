[![pipeline](https://github.com/jiekun/entry-task/actions/workflows/github-actions.yml/badge.svg)](https://github.com/jiekun/entry-task/actions)
[![codecov](https://codecov.io/gh/jiekun/entry-task/branch/master/graph/badge.svg?token=V0Y5Q4D3I0)](https://codecov.io/gh/jiekun/entry-task)
# Entry Task
## Introduction
`entry-task` is a web service based on **Gin** framework. It provides basic user management feature.

README is also available in Chinese: [中文文档](https://github.com/jiekun/entry-task/blob/master/README_CN.md).

**IMPORTANT**: This repo is set to public for study and discussion purposes. You may fork this project for reference usage. Please do **NOT** copy and paste directly when working on your own entry task.

## Document
You could find following guidance on [Entry Task System Design Document](https://docs.google.com/document/d/1sd5S8xdJRYcZrYAOM1cREnuQslZnkj8kIIQccHNmlq4/edit#) and `README.md`:
- Feature Description
- API Doc
- Deploy Doc
- Benchmark Result

## Setup
Before getting start, you need **MySQL** and **Redis** installed on your machine. And please refer `go.mod` for the Go packages used in this project. `init.sql` in script folder can help to initialize database & table we need.

Following commands will help you set up the project and run it locally:
```bash
# generate a config file
cp ./configs/config.yaml.defaul ./configs/config.yaml

# modify the host/db/etc settings in config file
vim ./configs/config.yaml

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

To run unit tests, try:
```bash
go test -cover -covermode=atomic -gcflags=all=-l ./... -coverprofile=profile.cov
go tool cover -func profile.cov
```

## Project Layout
```
├── LICENSE
├── README.md
├── cmd                                     # 
│    ├── grpc-server                        # 
│    │    └── main.go                       # gRPC Server entry point
│    └── http-server                        # 
│         └── main.go                       # HTTP Server entry point
├── configs                                 # 
│    └── config.yaml.default                # Config Template
├── global                                  # 
├── go.mod                                  # 
├── go.sum                                  # 
├── internal                                # Service Codes
│    ├── constant                           # 
│    ├── dao                                # DAO Layer
│    ├── error                              # 
│    ├── models                             # Models and DAO methods
│    ├── routers                            # 
│    │    ├── api                           # Controller Layer
│    │    └── router.go                     # Router for HTTP service
│    └── service                            # 
│         ├── grpc-service                  # Service Layer for gRPC service
│         └── http-service                  # Service Layer for HTTP service
├── log                                     # 
├── pkg                                     # Public utils
├── proto                                   # proto for RPC service
├── scripts                                 # Init & Test scripts 
└── upload                                  # 
```
