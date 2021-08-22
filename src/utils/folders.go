package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/user"
)

var LogFolder = getFolders()

func getFolders() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	return fmt.Sprint(usr.HomeDir, "/.kss/logs")
}
