// @Author: 2014BDuck
// @Date: 2021/5/16

package dao

import "github.com/2014bduck/entry-task/internal/models"

func (d *Dao) ValidateUser(username, password string) error {
	user := models.UserTab{Name: username, Password: password}
	_, err := user.GetValidUser(d.engine)
	if err != nil {
		return err
	}
	return nil
}
