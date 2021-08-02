// @Author: 2014BDuck
// @Date: 2021/8/2

package hashing

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/2014bduck/entry-task/global"
)

// HashPassword: Added salt and return a hashed string
func HashPassword(input string) string {
	hash := md5.Sum([]byte(global.AppSetting.HashSalt + input))
	return hex.EncodeToString(hash[:])
}
