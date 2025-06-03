package repositories

import (
	"auth-svc/pkg/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface{}

type UserRepository struct {
	BaseRepository[*models.User]
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{
		BaseRepository: BaseRepository[*models.User]{
			Collection: collection,
		},
	}
}
