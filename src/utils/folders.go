package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/user"
)

var AppFolder, LogFolder = getFolders()

func getFolders() (string, string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}

	appFolder := fmt.Sprint(usr.HomeDir, "/.kss")
	logFolder := fmt.Sprint(appFolder, "/logs")

	return appFolder, logFolder
}
