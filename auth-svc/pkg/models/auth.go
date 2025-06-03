package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	BaseModel `bson:",inline"`
	Email     string `json:"email" bson:"email,omitempty"`
	Password  string `json:"password" bson:"password,omitempty"`
}

func (u *User) SetID(id primitive.ObjectID) {
	u.ID = id
}

func (u *User) SetTimestamps() {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
}

type Admin struct {
	BaseModel
	Email    string `json:"email" bson:"email,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
}
