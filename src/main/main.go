package main

import (
	"context"
	"fmt"
	"github.com/ACLzz/keystore-server/src/api"
	"github.com/ACLzz/keystore-server/src/database"
	"github.com/ACLzz/keystore-server/src/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	log.Info("initializing keystore server")
	signal.Notify(utils.SigCh, os.Interrupt)
	utils.InitLogger()
	database.InitDb()
	
	s := startServer()
	finish(s)
}

func startServer() *http.Server {
	log.Info(fmt.Sprintf("Initializing http server on %s:%d", utils.Config.Addr, utils.Config.Port))
	s := http.Server{}
	
	mainR := api.MainRouter()
	s.Addr = fmt.Sprintf("%s:%d", utils.Config.Addr, utils.Config.Port)
	s.Handler = mainR
	go func() {
		log.Info("Starting keystore server...")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		utils.EndCh<-0
	}()
	
	return &s
}

func finish(s *http.Server) {
	for {
		select {
		case <-utils.SigCh:
			log.Info("Got interrupt signal. Exiting...")
			if err := s.Shutdown(context.Background()); err != nil {
				log.Fatalf("Error in shutting down server: %v\n", err)
			}
			os.Exit(0)

		case fin := <- utils.EndCh:
			log.Info("Shutting down server...")
			if err := s.Shutdown(context.Background()); err != nil {
				log.Fatalf("Error in shutting down server: %v\n", err)
			}
			os.Exit(fin)
		}
	}
}
