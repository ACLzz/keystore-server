package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type config struct {
	Addr 			string 	`yaml:"addr"`
	Port			int		`yaml:"port"`
	Dev				bool
	Test			bool
	
	DBHost			string	`yaml:"db_host"`
	DBPort			int		`yaml:"db_port"`
	DBName			string	`yaml:"db_name"`
	DBUsername		string	`yaml:"db_username"`
	DBPassword		string	`yaml:"db_password"`
	DSN				string

	Timezone		string	`yaml:"timezone"`
	TokenLifetime	int		`yaml:"token_lifetime"`
	Salt			string	`yaml:"salt"`
	AllwRegstr		bool	`yaml:"allow_registration"`
}

var Config = loadConfig()

func loadConfig() config {
	confObj := config{}
	confFn := "config.yml"

	confObj.Dev, confObj.Test = getMode()

	if confObj.Dev {
		confFn = fmt.Sprint("dev_", confFn)
	} else if confObj.Test{
		confFn = fmt.Sprint("test_", confFn)
	}

	path := fmt.Sprintf("%s/%s", AppFolder, confFn)
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Got error while reading %s config file: %v\n", path, err)
	}

	err = yaml.Unmarshal(f, &confObj)
	if err != nil {
		log.Fatalf("Got error while unmarshalling %s config file: %v\n", confFn, err)
	}
	
	if confObj.Test {
		confObj.DBName = "test_" + confObj.DBName
	}
	
	// Check timezone
	if _, err := time.LoadLocation(confObj.Timezone); err != nil {
		log.Fatal("Invalid timezone: ", confObj.Timezone)
	}

	confObj.DSN = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable TimeZone=%s",
		confObj.DBUsername, confObj.DBPassword, confObj.DBName, confObj.DBHost, confObj.DBPort, confObj.Timezone)
	return confObj
}


func getMode() (bool, bool) {
	dev, test := false, false
	if os.Getenv("MODE") == "dev" {
		dev = true
	} else if os.Getenv("MODE") == "test" {
		test = true
	}
	return dev, test
}
