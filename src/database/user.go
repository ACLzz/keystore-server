package database

import (
	"crypto/sha256"
	"fmt"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/ACLzz/keystore-server/src/utils"
	"time"
)

func (u *User) Register() error {
	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()

	if exists := u.isExist(u.Username); exists {
		return errors.UserExists
	}

	u.RegistrationDate = time.Now()

	// Hashing password
	saltyPassword := fmt.Sprint(u.Password, utils.Config.Salt)
	u.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(saltyPassword)))

	conn.Create(u)
	return nil
}

func (u *User) isExist(username string) bool {
	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()

	exists := User{}
	conn.Where("username = ?", username).First(&exists)

	if exists.Username != u.Username {
		return false
	}
	return true
}
