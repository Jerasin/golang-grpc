package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	DeletedAt time.Time          `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
