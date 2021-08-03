// @Author: 2014BDuck
// @Date: 2021/8/3

package main

import (
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/internal/models"
	"github.com/2014bduck/entry-task/internal/rpc-service"
	"github.com/2014bduck/entry-task/pkg/logger"
	"github.com/2014bduck/entry-task/pkg/rpc/tinyrpc"
	"github.com/2014bduck/entry-task/pkg/setting"
	rpcproto "github.com/2014bduck/entry-task/proto/rpc-proto"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
)

var (
	port    string
	runMode string
	config  string
)

// Server struct
type Server struct {
	addr  string
	funcs map[string]reflect.Value
}

// NewServer creates a new server
func NewServer(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

// Run server
func (s *Server) Run() {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("listen on %s err: %v\n", s.addr, err)
		return
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept err: %v\n", err)
			continue
		}

		go func() {
			srvTransport := tinyrpc.NewTransport(conn)

			for {
				// read request from client
				req, err := srvTransport.Receive()
				if err != nil {
					if err != io.EOF {
						log.Printf("read err: %v\n", err)
					}
					return
				}
				// get method by name
				f, ok := s.funcs[req.Name]
				if !ok { // if method requested does not exist
					e := fmt.Sprintf("func %s does not exist", req.Name)
					log.Println(e)
					if err = srvTransport.Send(tinyrpc.Data{Name: req.Name, Err: e}); err != nil {
						log.Printf("transport write err: %v\n", err)
					}
					continue
				}
				log.Printf("func %s is called\n", req.Name)
				// unpack request arguments
				inArgs := make([]reflect.Value, len(req.Args))
				for i := range req.Args {
					inArgs[i] = reflect.ValueOf(req.Args[i])
				}
				// invoke requested method
				out := f.Call(inArgs)
				// pack response arguments (except error)
				outArgs := make([]interface{}, len(out)-1)
				for i := 0; i < len(out)-1; i++ {
					outArgs[i] = out[i].Interface()
				}
				// pack error argument
				var e string
				if _, ok := out[len(out)-1].Interface().(error); !ok {
					e = ""
				} else {
					e = out[len(out)-1].Interface().(error).Error()
				}
				// send response to client
				err = srvTransport.Send(tinyrpc.Data{Name: req.Name, Args: outArgs, Err: e})
				if err != nil {
					log.Printf("transport write err: %v\n", err)
				}
			}
		}()
	}
}

// Register a method via name
func (s *Server) Register(name string, f interface{}) {
	if _, ok := s.funcs[name]; ok {
		return
	}
	s.funcs[name] = reflect.ValueOf(f)
}

func main() {
	s := NewServer(":" + global.ServerSetting.RPCPort)
	c := context.Background()
	svc := rpcservice.New(c)
	registerRPC(s, svc)
	log.Printf("Starting RPC server, Listening %s...\n", port)
	go s.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	log.Println("Server exited")
}

func init() {
	err := setupFlag()
	if err != nil {
		log.Fatalf("init.setupFlag err: %v", err)
	}

	err = setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	err = setupCacheClient()
	if err != nil {
		log.Fatalf("init.setupCacheClient err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
	flag.Parse()
	return nil
}

func setupSetting() error {
	log.Printf("%v", config)
	s, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}
	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("Redis", &global.CacheSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("Client", &global.ClientSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	if port != "" {
		global.ServerSetting.RPCPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}
	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = models.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupCacheClient() error {
	var err error
	global.CacheClient, err = models.NewCacheClient(global.CacheSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(
		&lumberjack.Logger{
			Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
			MaxSize:   600,
			MaxAge:    10,
			LocalTime: true,
		},
		"",
		log.LstdFlags,
	).WithCaller(2)

	return nil
}

func registerRPC(s *Server, svc rpcservice.Service) {
	gob.Register(rpcproto.UserLoginRequest{})
	gob.Register(rpcproto.UserLoginResponse{})
	gob.Register(rpcproto.UserRegisterRequest{})
	gob.Register(rpcproto.UserRegisterResponse{})
	gob.Register(rpcproto.UserEditRequest{})
	gob.Register(rpcproto.UserEditResponse{})
	gob.Register(rpcproto.UserGetRequest{})
	gob.Register(rpcproto.UserGetResponse{})
	s.Register("Register", svc.UserRegister)
	s.Register("Login", svc.UserLogin)
	s.Register("EditUser", svc.UserEdit)
	s.Register("GetUser", svc.UserGet)
}
