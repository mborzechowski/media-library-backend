package routes

import (
	"game-catalog-backend/controllers"

	"github.com/gorilla/mux"
)

func RegisterGameRoutes(router *mux.Router) {
	router.HandleFunc("/games", controllers.GetGames).Methods("GET")
	router.HandleFunc("/games", controllers.AddGame).Methods("POST")
	router.HandleFunc("/games", controllers.HandleOptions).Methods("OPTIONS")
	router.HandleFunc("/games/{id}", controllers.GetGame).Methods("GET")

}
