// @Author: 2014BDuck
// @Date: 2021/7/11

package global

import (
	"github.com/jiekun/entry-task/pkg/logger"
	"github.com/jiekun/entry-task/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	CacheSetting    *setting.CacheSettingS
	ClientSetting   *setting.ClientSettingS
	Logger          *logger.Logger
)
