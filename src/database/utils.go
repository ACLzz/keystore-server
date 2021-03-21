package database

import (
	"github.com/ACLzz/keystore-server/src/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetConn() *gorm.DB {
	db, err := gorm.Open(postgres.Open(utils.Config.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("Can't get connection: ", err)
	}
	return db
}


func InitDb() {
	log.Info("Initializing database...")
	db := GetConn()

	DB, _ := db.DB()
	defer DB.Close()
	
	for _, model := range ModelsList {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatal(err)
		}
	}
}
