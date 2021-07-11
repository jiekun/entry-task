// @Author: 2014BDuck
// @Date: 2021/7/11

package global

import (
	"github.com/2014bduck/entry-task/pkg/logger"
	"github.com/2014bduck/entry-task/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
)
