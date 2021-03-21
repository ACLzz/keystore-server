package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type config struct {
	Addr 		string 	`yaml:"addr"`
	Port		int 	`yaml:"port"`
	
	DBHost		string	`yaml:"db_host"`
	DBPort		int		`yaml:"db_port"`
	DBName		string	`yaml:"db_name"`
	DBUsername	string	`yaml:"db_username"`
	DBPassword	string	`yaml:"db_password"`
	DSN			string
}

var Config = loadConfig()
var Dev, Test = getMode()

func loadConfig() config {
	confObj := config{}
	confFn := "config.yml"

	if Dev {
		confFn = fmt.Sprint("dev_", confFn)
	} else if Test{
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
	fmt.Println(confObj.DBUsername)

	confObj.DSN = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable TimeZone=%s",
		confObj.DBUsername, confObj.DBPassword, confObj.DBName, confObj.DBHost, confObj.DBPort, "Europe/Zaporozhye") // FIXME
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
