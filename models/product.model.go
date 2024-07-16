package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Image       string             `json:"image"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Price       int                `json:"price"`
}
