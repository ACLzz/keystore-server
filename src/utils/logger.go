package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

func InitLogger() {
	logPath := fmt.Sprint(LogFolder, "/", time.Now().Format("02-01-2006_15:04"), ".log")
	if f, err := os.OpenFile(logPath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0755); err != nil {
		log.Error("Open file for write log: ", err)
		return
	} else {
		log.SetOutput(io.MultiWriter(os.Stderr, f))
	}
}
