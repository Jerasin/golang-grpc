package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepository[T any] interface {
	Save(data T) (T, error)
	IsExist(filters primitive.M) (T, error)
	FindByID(id primitive.ObjectID) (T, error)
	FindAll(filters map[string]any) ([]T, error)
	Update(id primitive.ObjectID, data T) (T, error)
	Delete(id primitive.ObjectID) error
	Count(filters map[string]any) (int64, error)
}

type Timestamped interface {
	SetID(primitive.ObjectID)
	SetTimestamps()
}

type baseRepository[T Timestamped] struct {
	Collection *mongo.Collection
}

func NewBaseRepository[T Timestamped](collection *mongo.Collection) BaseRepository[T] {
	return &baseRepository[T]{Collection: collection}
}

func MakeFilter(filters map[string]any) primitive.M {
	filter := primitive.M{"deleted_at": nil} // Default filter to exclude soft-deleted records
	for key, value := range filters {
		if value != nil {
			filter[key] = value
		}
	}
	return filter
}

func (r *baseRepository[T]) Save(data T) (T, error) {
	data.SetID(primitive.NewObjectID())
	data.SetTimestamps()

	_, err := r.Collection.InsertOne(context.Background(), data)
	if err != nil {
		var zero T
		return zero, err
	}
	return data, nil
}

func (r *baseRepository[T]) IsExist(filters primitive.M) (T, error) {
	var zero T
	result := r.Collection.FindOne(context.Background(), filters)

	if result.Err() != nil {
		return zero, result.Err()
	}

	var data T
	err := result.Decode(&data)
	if err != nil {
		return zero, err
	}

	data.SetID(primitive.NewObjectID())
	return data, nil
}

func (r *baseRepository[T]) FindByID(id primitive.ObjectID) (T, error) {
	var zero T
	filter := MakeFilter(map[string]any{"_id": id})
	result := r.Collection.FindOne(context.Background(), filter)

	if result.Err() != nil {
		return zero, result.Err()
	}

	var data T
	err := result.Decode(&data)
	if err != nil {
		return zero, err
	}

	data.SetID(id)
	return data, nil
}
func (r *baseRepository[T]) FindAll(filters map[string]any) ([]T, error) {
	filter := MakeFilter(filters)
	cursor, err := r.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []T
	for cursor.Next(context.Background()) {
		var data T
		err := cursor.Decode(&data)
		if err != nil {
			return nil, err
		}
		data.SetID(primitive.NewObjectID())
		results = append(results, data)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *baseRepository[T]) Update(id primitive.ObjectID, data T) (T, error) {
	data.SetID(id)
	data.SetTimestamps()

	filter := MakeFilter(map[string]any{"_id": id})
	update := primitive.M{"$set": data}

	_, err := r.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		var zero T
		return zero, err
	}

	return data, nil
}

func (r *baseRepository[T]) Delete(id primitive.ObjectID) error {
	filter := MakeFilter(map[string]any{"_id": id})
	_, err := r.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *baseRepository[T]) Count(filters map[string]any) (int64, error) {
	filter := MakeFilter(filters)
	count, err := r.Collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
