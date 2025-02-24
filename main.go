package main

import (
	"log"
	"net/http"
	
	"kubecodya/config"
	"kubecodya/middleware"
	"kubecodya/routes"
)

func main() {
	router := routes.NewRouter()
	loggedRouter := middleware.Logging(router)
	
	log.Printf("Server running on port %s", config.Port)
	if err := http.ListenAndServe(config.Port, loggedRouter); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
