package main

import (
	"context"
	"fmt"
	"github.com/ACLzz/AndroidBeaconService-server/src/api"
	"github.com/ACLzz/AndroidBeaconService-server/src/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	log.Info("initializing keystore server")
	signal.Notify(utils.SigCh, os.Interrupt)
	utils.InitLogger()
	
	s := startServer()
	finish(s)
}

func startServer() *http.Server {
	log.Info("Initializing http server")
	s := http.Server{}
	c := &utils.Config
	
	mainR := api.MainRouter()
	s.Addr = fmt.Sprintf("%s:%d", c.Addr, c.Port)
	s.Handler = mainR
	go func() {
		log.Info("Starting keystore server...")
		if err := s.ListenAndServe(); err != nil {
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
