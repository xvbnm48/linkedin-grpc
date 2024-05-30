package main

import (
	"github.com/xvbnm48/linkedin-grpc/internal/database"
	"github.com/xvbnm48/linkedin-grpc/internal/server"
	"log"
)

func main() {
	db, err := database.NewDatabaseClient()
	if err != nil {
		log.Fatalf("failed to initialize: %v", err)
	}
	srv := server.NewEchoServer(db)
	if srv := srv.Start(); srv != nil {
		log.Fatalf(err.Error())
	}

}
