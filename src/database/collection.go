package database

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (c *Collection) Add() bool {
	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()

	conn.Create(c)
	return true
}

func (c *Collection) IsExist() bool {
	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()

	exists := Collection{}
	if tx := conn.Where("title = ? AND user_refer = ?", c.Title, c.User.Id).First(&exists); tx.Error == gorm.ErrRecordNotFound {
		return false
	}
	c.Id = exists.Id
	return true
}

func (c *Collection) Update() bool {
	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	defer conn.Commit()

	if tx := conn.Unscoped().Save(c); tx.Error != nil {
		log.Error(tx.Error)
		return false
	}
	return true
}

func (c *Collection) Delete() bool {
	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	defer conn.Commit()
	
	if tx := conn.Unscoped().Where("id = ?", c.Id).Delete(Collection{}); tx.Error != nil {
		log.Error(tx.Error)
		return false
	}
	return true
}

func (c *Collection) FetchPasswords() []Password {
	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	var passwords []Password

	if tx := conn.Where("collection_refer = ?", c.Id).Find(&passwords); tx.Error != nil {
		log.Error(tx.Error)
		return nil
	}
	return passwords
}