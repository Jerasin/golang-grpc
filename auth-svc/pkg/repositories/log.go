package repositories

import (
	"auth-svc/pkg/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuditLogRepositoryInterface interface{}

type AuditLogRepository struct {
	BaseRepository[*models.AuditLog]
}

func NewAuditLogRepository(collection *mongo.Collection) *AuditLogRepository {
	return &AuditLogRepository{
		BaseRepository: NewBaseRepository[*models.AuditLog](collection),
	}
}
