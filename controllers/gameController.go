package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"game-catalog-backend/config"
	"game-catalog-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetGames retrieves all games from the database
func GetGames(w http.ResponseWriter, r *http.Request) {
	var games []models.Game
	collection := config.Client.Database("gamecatalog").Collection("games")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var game models.Game
		cursor.Decode(&game)
		games = append(games, game)
	}

	json.NewEncoder(w).Encode(games)
}

// AddGame adds a new game to the database
func AddGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	json.NewDecoder(r.Body).Decode(&game)

	collection := config.Client.Database("gamecatalog").Collection("games")
	result, err := collection.InsertOne(context.Background(), game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the inserted ID and set it in the response
	id := result.InsertedID.(primitive.ObjectID)
	game.ID = id

	log.Printf("Received game: %+v\n", game)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(game)
}

// Pobieranie pojedynczej gry
func GetGame(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var game models.Game
    collection := config.Client.Database("gamecatalog").Collection("games")
    filter := bson.M{"_id": objID}

    err = collection.FindOne(context.Background(), filter).Decode(&game)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(game)
}


// HandleOptions handles CORS preflight requests
func HandleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}
