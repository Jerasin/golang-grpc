package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Timestamped interface {
	SetID(primitive.ObjectID)
	SetTimestamps()
}

type BaseRepository[T Timestamped] struct {
	Collection *mongo.Collection
}

func (r *BaseRepository[T]) Register(data T) (T, error) {
	data.SetID(primitive.NewObjectID())
	data.SetTimestamps()

	_, err := r.Collection.InsertOne(context.Background(), data)
	if err != nil {
		var zero T
		return zero, err
	}
	return data, nil
}

func (r *BaseRepository[T]) IsExist(filters primitive.M) (T, error) {
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

func (r *BaseRepository[T]) FindByID(id primitive.ObjectID) (T, error) {
	var zero T
	filter := primitive.M{"_id": id}
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
func (r *BaseRepository[T]) FindAll() ([]T, error) {
	filter := primitive.M{}
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

func (r *BaseRepository[T]) Update(id primitive.ObjectID, data T) (T, error) {
	data.SetID(id)
	data.SetTimestamps()

	filter := primitive.M{"_id": id}
	update := primitive.M{"$set": data}

	_, err := r.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		var zero T
		return zero, err
	}

	return data, nil
}

func (r *BaseRepository[T]) Delete(id primitive.ObjectID) error {
	filter := primitive.M{"_id": id}
	_, err := r.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
