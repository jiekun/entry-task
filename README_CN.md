[![pipeline](https://github.com/jiekun/entry-task/actions/workflows/github-actions.yml/badge.svg)](https://github.com/jiekun/entry-task/actions)
[![codecov](https://codecov.io/gh/jiekun/entry-task/branch/master/graph/badge.svg?token=V0Y5Q4D3I0)](https://codecov.io/gh/jiekun/entry-task)
# Entry Task
## 简介
`entry-task` 是一个基于 **Gin** 框架的 web service，提供基础的用户管理功能。

## 文档
你可以从 [Entry Task 系统设计说明书](https://docs.google.com/document/d/1sd5S8xdJRYcZrYAOM1cREnuQslZnkj8kIIQccHNmlq4/edit#) 和 `README.md` 中查阅以下内容:
- 功能描述
- API 文档
- Deploy 部署文档
- 性能测试结果

## 配置运行
在启动之前, 你需要在本地安装 **MySQL** and **Redis**。 请从 `go.mod` 中查看本项目使用到的第三方库。 script 目录中的 `init.sql` 文件包含了初始化数据库所需要的语句。

以下命令能帮助你在本地配置项目并运行:
```bash
# 生成一份配置文件
cp ./configs/config.yaml.defaul ./configs/config.yaml

# 修改配置文件中的 host/db 及其他配置
vim ./configs/config.yaml

# 处理第三方库依赖
go mod tidy

# 启动 RPC server
go run ./cmd/grpc-server/main.go

# 启动 HTTP server
go run ./cmd/http-server/main.go 
```

最后测试一下是否正确运行:
```bash
curl --location --request GET 'http://{your_http_host}/api/ping' # {"message":"pong"}
```

用以下命令执行单元测试:
```bash
go test -cover -covermode=atomic -gcflags=all=-l ./... -coverprofile=profile.cov
go tool cover -func profile.cov
```

## 项目结构
```
├── LICENSE
├── README.md
├── cmd                                     # 
│    ├── rpc-server                         # 
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
