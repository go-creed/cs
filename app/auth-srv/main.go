package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"cs/app/auth-srv/handler"
	"cs/app/auth-srv/subscriber"

	auth "cs/app/auth-srv/proto/auth"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.cs.service.auth"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	auth.RegisterAuthHandler(service.Server(), new(handler.Auth))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.cs.service.auth", service.Server(), new(subscriber.Auth))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
