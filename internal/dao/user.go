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
		Status:     status,
		CommonModel: &models.CommonModel{
			Ctime: uint32(time.Now().Unix()),
			Mtime: uint32(time.Now().Unix()),
		},
	}
	return userTab.Create(d.engine)

}

func (d *Dao) UpdateUser(id uint32, nickname, profilePic string) error {
	user := models.UserTab{CommonModel: &models.CommonModel{ID: id}}
	values := map[string]interface{}{
		"mtime": uint32(time.Now().Unix()),
	}
	if profilePic != ""{
		values["profile_pic"] = profilePic
	}
	if nickname != ""{
		values["nickname"] = nickname
	}

	return user.Update(d.engine, values)

}
