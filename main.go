package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func main() {
	config := getConfig()

	server := NewServer(config)
	server.Initialize()

	defer func() {
		if err := server.mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	server.Listen()

	<-server.connClose
	log.Info("Shutdown complete")
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
}
