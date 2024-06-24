package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Game struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Genre           string   `json:"genre"`
	Platform        string   `json:"platform"`
	YearPublished   int      `json:"yearPublished"`
	PhysicalDigital string   `json:"physicalDigital"`
	Publisher       string   `json:"publisher"`
	Developer       string   `json:"developer"`
	Completed       bool     `json:"completed"`
	Rating          int      `json:"rating"`
	Images          []string `json:"images"`
}
func init() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func connectDB() {
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        log.Fatal("MONGODB_URI is not set")
    }
    clientOptions := options.Client().ApplyURI(uri)
    var err error
    client, err = mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Test connection
    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal("Cannot connect to MongoDB:", err)
    }
    log.Println("Connected to MongoDB")
}

func getGames(w http.ResponseWriter, r *http.Request) {
    var games []Game
    collection := client.Database("gamecatalog").Collection("games")
    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())
    for cursor.Next(context.Background()) {
        var game Game
        cursor.Decode(&game)
        games = append(games, game)
    }
    json.NewEncoder(w).Encode(games)
}

func addGame(w http.ResponseWriter, r *http.Request) {
    var game Game
    json.NewDecoder(r.Body).Decode(&game)
    collection := client.Database("gamecatalog").Collection("games")
    _, err := collection.InsertOne(context.Background(), game)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Received game: %+v\n", game)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(game)
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusOK)
}

func main() {
    connectDB()
    router := mux.NewRouter()

    corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"http://localhost:3000"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
        handlers.AllowedHeaders([]string{"Content-Type"}),
    )

    router.Use(corsHandler)

    router.HandleFunc("/games", getGames).Methods("GET")
    router.HandleFunc("/games", addGame).Methods("POST")
    router.HandleFunc("/games", handleOptions).Methods("OPTIONS")

    log.Fatal(http.ListenAndServe(":8080", router))
}