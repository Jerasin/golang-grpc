package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLog struct {
	BaseModel `bson:",inline"`
	Request   string   `json:"request" bson:"request,omitempty"`
	Response  string   `json:"response" bson:"response,omitempty"`
	UserID    string   `json:"user_id" bson:"user_id,omitempty"`
	Method    string   `json:"method" bson:"method,omitempty"`
	ClientIP  string   `json:"client_ip" bson:"client_ip,omitempty"`
	Status    string   `json:"status" bson:"status,omitempty"`
	Error     string   `json:"error" bson:"error,omitempty"`
	Duration  int64    `json:"duration" bson:"duration,omitempty"`
	Metadata  string   `json:"meta_data" bson:"meta_data,omitempty"`
	Protocol  string   `json:"protocol" bson:"protocol,omitempty"`
	TraceID   []string `json:"trace_id" bson:"trace_id,omitempty"`
}

func (a *AuditLog) SetID(id primitive.ObjectID) {
	a.ID = id
}

func (a *AuditLog) GetID() primitive.ObjectID {
	return a.ID
}

func (a *AuditLog) SetTimestamps() {
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
}
