package main

import (
        log "github.com/micro/go-micro/v2/logger"
        "github.com/micro/go-micro/v2/web"
)

func main() {
	// create new web service
        service := web.NewService(
                web.Name("go.micro.cs.web.user"),
                web.Version("latest"),
        )

	// initialise service
        if err := service.Init(); err != nil {
                log.Fatal(err)
        }

	// register html handler
	//service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	//service.HandleFunc("/user/call", handler.UserCall)

	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
