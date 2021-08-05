package main

import (
	"context"
	"errors"
	"flag"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/internal/models"
	"github.com/2014bduck/entry-task/internal/routers"
	"github.com/2014bduck/entry-task/internal/service"
	"github.com/2014bduck/entry-task/pkg/logger"
	"github.com/2014bduck/entry-task/pkg/rpc/erpc"
	"github.com/2014bduck/entry-task/pkg/rpc/grpc"
	"github.com/2014bduck/entry-task/pkg/setting"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	port    string
	runMode string
	config  string
)

func main() {
	r := routers.NewRouter()

	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        r,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("Starting HTTP server, Listening %s...\n", s.Addr)
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("s.ListenAndServe err: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
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

	err = setupRPCClient()
	if err != nil {
		log.Fatalf("init.setupRPCClient err: %v", err)
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
		global.ServerSetting.HttpPort = port
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

func setupRPCClient() error {
	ctx := context.Background()
	clientConn, err := grpc.GetClientConn(ctx, global.ClientSetting.RPCHost, nil)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	global.GRPCClient = clientConn

	conn, err := net.Dial("tcp", global.ClientSetting.RPCHost)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	global.RPCClient = erpc.NewClient(conn)
	service.RegisterUserServiceProto()
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
