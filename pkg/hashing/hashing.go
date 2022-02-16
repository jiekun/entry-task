// @Author: 2014BDuck
// @Date: 2021/8/2

package hashing

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/jiekun/entry-task/global"
	"golang.org/x/crypto/bcrypt"
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

func HashPasswordBcrypt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHashBcrypt(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
