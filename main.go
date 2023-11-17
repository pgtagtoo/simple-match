package main

import (
	"log"
	"net/http"

	"github.com/genepg/simple-match-api/api"
)

func main() {
	s := api.CreateNewServer()
	s.MountHandlers()

	port := "8080"
	log.Println("Server listen on " + port)
	http.ListenAndServe(":"+port, s.Router)
}
