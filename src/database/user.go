package database

import (
	"crypto/sha256"
	"fmt"
	"github.com/ACLzz/keystore-server/src/errors"
	"github.com/ACLzz/keystore-server/src/utils"
	log "github.com/sirupsen/logrus"
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

func (u *User) Delete() bool {
	var _u User
	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	defer conn.Commit()

	conn.Unscoped().Where("username = ?", u.Username).First(&_u)
	if tx := conn.Unscoped().Delete(&_u); tx.Error != nil {
		log.Error(tx.Error)
		return false
	}

	if tx := conn.Unscoped().Where("user_refer = ?", _u.Id).Delete(Token{}); tx.Error != nil {
		log.Error(tx.Error)
		return false
	}

	if tx := conn.Unscoped().Where("username = ?", u.Username).Delete(User{}); tx.Error != nil {
		log.Error(tx.Error)
		return false
	}
	return true
}

func (u *User) Update() bool {
	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	defer conn.Commit()
	
	if tx := conn.Unscoped().Save(u); tx.Error != nil {
		log.Error(tx.Error)
		return false
	}
	return true
}
