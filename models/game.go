package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Game struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title           string             `json:"title"`
	Genre           string             `json:"genre"`
	Platform        string             `json:"platform"`
	YearPublished   int                `json:"yearPublished"`
	PhysicalDigital string             `json:"physicalDigital"`
	Publisher       string             `json:"publisher"`
	Developer       string             `json:"developer"`
	Completed       bool               `json:"completed"`
	Rating          int                `json:"rating"`
	Images          []string           `json:"images"`
}
