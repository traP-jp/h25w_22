package main

import (
	"log"
	"net/http"

	"github.com/traP-jp/h25w_22/handler"
)

func main() {
	manager := handler.NewRoomManager()
	mux := handler.SetupRoutes(manager)

	port := ":8080"
	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
