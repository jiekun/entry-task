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
	return HashWithMD5(global.AppSetting.HashSalt + input)
}

// HashWithMD5: Return a MD5'd string
func HashWithMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
