package database

import (
	log "github.com/sirupsen/logrus"
)

func (p *Password) Add(collection string, userId int) bool {
	conn := GetConn()
	DB, _ := conn.DB()
	defer DB.Close()
	
	coll := Collection{}
	if tx := conn.Where("title = ? and user_refer = ?", collection, userId).First(&coll); tx.Error != nil {
		log.Error(tx.Error)
		return false
	}
	
	p.Collection = coll
	if tx := conn.Create(p); tx.Error != nil {
		log.Error(tx.Error)
		return false
	}
	return true
}
