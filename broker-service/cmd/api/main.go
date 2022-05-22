package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = 80

type Config struct {
}

func main() {
	app := Config{}

	log.Printf("starting broker service on port %d\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
