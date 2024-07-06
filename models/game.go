package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Game struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title           string             `json:"title" bson:"title"`
	Genre           string             `json:"genre" bson:"genre"`
	Platform        string             `json:"platform" bson:"platform"`
	YearPublished   int                `json:"yearPublished" bson:"yearPublished"`
	PhysicalDigital string             `json:"physicalDigital" bson:"physicalDigital"`
	Publisher       string             `json:"publisher" bson:"publisher"`
	Developer       string             `json:"developer" bson:"developer"`
	Completed       bool               `json:"completed" bson:"completed"`
	Rating          int                `json:"rating" bson:"rating"`
	Images          []string           `json:"images" bson:"images"`
}
