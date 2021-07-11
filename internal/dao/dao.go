// @Author: 2014BDuck
// @Date: 2021/7/11

package dao

import "gorm.io/gorm"

type Dao struct {
	engine *gorm.DB
}

func New(engine *gorm.DB) *Dao {
	return &Dao{engine: engine}
}
