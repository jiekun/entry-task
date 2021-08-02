// @Author: 2014BDuck
// @Date: 2021/5/16

package dao

import (
	"github.com/2014bduck/entry-task/internal/models"
	"time"
)

func (d *Dao) GetUserByName(username string) (models.UserTab, error) {
	userTab := models.UserTab{Name: username}
	user, err := userTab.GetUserByName(d.engine)

	// ErrRecordNotFound or DB error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (d *Dao) CreateUser(username, password, nickname, profilePic string, status uint8) (*models.UserTab, error) {
	userTab := models.UserTab{
		Name:       username,
		Password:   password,
		Nickname:   nickname,
		ProfilePic: profilePic,
		Status: status,
		CommonModel: &models.CommonModel{
			Ctime: uint32(time.Now().Unix()),
			Mtime: uint32(time.Now().Unix()),
		},
	}
	return userTab.Create(d.engine)

}
