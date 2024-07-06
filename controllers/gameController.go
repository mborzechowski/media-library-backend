package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"game-catalog-backend/config"
	"game-catalog-backend/models"
	"log"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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

    // Parse the incoming request
    err := r.ParseMultipartForm(10 << 20) // 10MB
    if err != nil {
        http.Error(w, "Unable to parse form data", http.StatusBadRequest)
        return
    }

    log.Println("Form data parsed successfully")

    // Retrieve file
    file, _, err := r.FormFile("image")
    if err != nil {
        log.Printf("Unable to read image from request: %v", err)
        http.Error(w, "Unable to read image from request", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Upload file to Cloudinary
    uploadResult, err := config.Cloudinary.Upload.Upload(context.Background(), file, uploader.UploadParams{})
    if err != nil {
        log.Printf("Failed to upload image: %v", err)
        http.Error(w, fmt.Sprintf("Failed to upload image: %v", err), http.StatusInternalServerError)
        return
    }

   
    // Decode other fields
    game.Title = r.FormValue("title")
    game.Genre = r.FormValue("genre")
    game.Platform = r.FormValue("platform")
    game.YearPublished, err = strconv.Atoi(r.FormValue("yearPublished"))
    if err != nil {
        log.Printf("Error parsing YearPublished: %v", err)
    }
    game.PhysicalDigital = r.FormValue("physicalDigital")
    game.Publisher = r.FormValue("publisher")
    game.Developer = r.FormValue("developer")
    game.Completed, err = strconv.ParseBool(r.FormValue("completed"))
    if err != nil {
        log.Printf("Error parsing Completed: %v", err)
    }
    game.Rating, err = strconv.Atoi(r.FormValue("rating"))
    if err != nil {
        log.Printf("Error parsing Rating: %v", err)
    }
    game.Images = []string{uploadResult.SecureURL}

    log.Printf("Game object populated: %+v\n", game)

    collection := config.Client.Database("gamecatalog").Collection("games")
    result, err := collection.InsertOne(context.Background(), game)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    id := result.InsertedID.(primitive.ObjectID)
    game.ID = id

    log.Printf("Game inserted into database with ID: %s\n", id.Hex())

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
