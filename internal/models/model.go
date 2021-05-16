// @Author: 2014BDuck
// @Date: 2021/5/16

package models

type CommonModel struct {
	ID    uint32 `gorm:"primary_key" json:"id,omitempty"`
	Ctime uint32 `json:"ctime,omitempty"`
	Mtime uint32 `json:"mtime,omitempty"`
}