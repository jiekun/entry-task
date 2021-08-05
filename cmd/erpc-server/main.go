// @Author: 2014BDuck
// @Date: 2021/8/3

package main

import (
	"context"
	"flag"
	"github.com/2014bduck/entry-task/global"
	service "github.com/2014bduck/entry-task/internal/erpc-service"
	"github.com/2014bduck/entry-task/internal/models"
	"github.com/2014bduck/entry-task/pkg/logger"
	"github.com/2014bduck/entry-task/pkg/rpc/erpc"
	"github.com/2014bduck/entry-task/pkg/setting"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
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

func main() {
	s := erpc.NewServer(":" + global.ServerSetting.RPCPort)
	c := context.Background()

	userSvc := service.NewUserService(c)
	userSvc.RegisterUserService(s)
	uploadSvc := service.NewUploadService(c)
	uploadSvc.RegisterUploadService(s)

	log.Printf("eRPC server listening: %v\n", global.ServerSetting.RPCPort)
	s.Run()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	log.Println("Server exited")
}
