// @Author: 2014BDuck
// @Date: 2021/5/16

package models

import (
	"fmt"
	"github.com/jiekun/entry-task/global"
	"github.com/jiekun/entry-task/pkg/setting"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type CommonModel struct {
	ID    uint32 `gorm:"primary_key" json:"id,omitempty"`
	Ctime uint32 `json:"ctime,omitempty"`
	Mtime uint32 `json:"mtime,omitempty"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	)))
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,       // Disable color
			},
		)
		db.Logger = newLogger
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}

func NewCacheClient(cacheSetting *setting.CacheSettingS) (*redis.Client, error) {
	rClient := redis.NewClient(&redis.Options{
		Addr: cacheSetting.Host,
		DB:   cacheSetting.DBIndex, // use default DB
	})
	return rClient, nil
}
