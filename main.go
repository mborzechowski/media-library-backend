package main

import (
	"game-catalog-backend/config"
	"game-catalog-backend/routes"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()
	router := mux.NewRouter()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	router.Use(corsHandler)

	routes.RegisterGameRoutes(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}