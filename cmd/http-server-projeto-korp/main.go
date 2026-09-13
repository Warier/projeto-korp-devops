package main

import (
	"log"

	"http-server-projeto-korp/internal/httpapi"
)

func main() {
	router, err := httpapi.NewRouter()
	if err != nil {
		log.Fatal(err)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
