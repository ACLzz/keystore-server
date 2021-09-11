package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	Addr 			string
	Port			int
	Dev				bool
	Test			bool
	
	DBHost			string
	DBPort			int
	DBName			string
	DBUsername		string
	DBPassword		string
	DSN				string

	Timezone		string
	TokenLifetime	int
	Salt			string
	AllwRegstr		bool
}

var Config = loadConfig()

func loadConfig() config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Warning(".env file doesn't located in directory, default values are set")
	}
	getStr := func(key string, def string) string {
		val, ok := os.LookupEnv(key)
		if !ok {
			log.Warningf("variable %s is not set, setting to default", key)
			return def
		}
		return val
	}
	getInt := func(key string, def int) int {
		val, ok := os.LookupEnv(key)
		if !ok {
			log.Warningf("variable %s is not set, setting to default", key)
			return def
		}
		i, _ := strconv.Atoi(val)
		return i
	}

	confObj := config{
		Addr: getStr("ADDR", "127.0.0.1"),
		Port: getInt("PORT", 8402),

		DBHost: getStr("POSTGRES_HOST", "localhost"),
		DBPort: getInt("POSTGRES_PORT", 5432),
		DBName: getStr("POSTGRES_NAME", "keystore"),
		DBUsername: getStr("POSTGRES_USERNAME", "keykeeper"),
		DBPassword: getStr("POSTGRES_PASSWORD", "CHANGE_PASSWORD"),

		Timezone: getStr("TIMEZONE", "Europe/Zaporozhye"),
		TokenLifetime: getInt("TOKEN_LIFETIME", 3600),
		Salt: getStr("SALT", "CHANGE_SALT"),
		AllwRegstr: strings.ToLower(getStr("ALLOW_REGISTRATION", "true")) == "true",
		Dev: strings.ToLower(getStr("MODE", "dev")) == "dev",
		Test: strings.ToLower(getStr("MODE", "test")) == "test",
	}
	confObj.DSN = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable TimeZone=%s",
		confObj.DBUsername, confObj.DBPassword, confObj.DBName, confObj.DBHost, confObj.DBPort, confObj.Timezone)
	
	// Check timezone
	if _, err := time.LoadLocation(confObj.Timezone); err != nil {
		log.Fatal("Invalid timezone: ", confObj.Timezone)
	}
	return confObj
}
