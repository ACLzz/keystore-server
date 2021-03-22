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

	if u.IsExist() {
		return errors.UserExists
	}

	u.RegistrationDate = time.Now()

	// Hashing password
	saltyPassword := fmt.Sprint(u.Password, utils.Config.Salt)
	u.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(saltyPassword)))

	conn.Create(u)
	return nil
}

func (u *User) IsExist() bool {
	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()

	exists := User{}
	conn.Where("username = ?", u.Username).First(&exists)

	if exists.Username != u.Username {
		return false
	}
	return true
}

func (u *User) GenToken() string {
	t := Token{
		User:         *u,
	}
	t.genToken()

	return t.Token
}

func (u *User) CheckPassword() bool {
	password := u.HashPassword()

	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()

	exists := User{}
	conn.Where("username = ? AND password = ?", u.Username, password).First(&exists)

	if exists.Username != u.Username {
		return false
	}
	return true
}

func (u *User) HashPassword() string {
	saltyPassword := fmt.Sprint(u.Password, utils.Config.Salt)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(saltyPassword)))
}
