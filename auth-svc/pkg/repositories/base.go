package repositories

import (
	"context"

	"github.com/goforj/godump"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BaseRepository[T any] interface {
	Save(ctx context.Context, data T) (T, error)
	IsExist(ctx context.Context, filters primitive.M) (T, error)
	FindOne(ctx context.Context, filters map[string]any) (T, error)
	FindAll(ctx context.Context, filters map[string]any) ([]T, error)
	Update(ctx context.Context, id primitive.ObjectID, data T) (T, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	Count(ctx context.Context, filters map[string]any) (int64, error)
	Transaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
}

type Timestamped interface {
	GetID() primitive.ObjectID
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

func (r *baseRepository[T]) Save(ctx context.Context, data T) (T, error) {
	data.SetID(primitive.NewObjectID())
	data.SetTimestamps()

	_, err := r.Collection.InsertOne(ctx, data)
	if err != nil {
		var zero T
		return zero, err
	}

	return data, nil
}

func (r *baseRepository[T]) IsExist(ctx context.Context, filters primitive.M) (T, error) {
	var resource T
	result := r.Collection.FindOne(ctx, filters)

	godump.Dump("Err", result.Err())

	if result.Err() != nil {

		if result.Err() == mongo.ErrNoDocuments {
			return resource, nil // No document found, return zero value
		}

		return resource, result.Err()
	}

	var data T
	err := result.Decode(&data)
	if err != nil {
		return resource, err
	}

	godump.Dump("check", data.GetID() == (primitive.NilObjectID))
	if !(data.GetID() == (primitive.NilObjectID)) {
		return resource, status.Error(codes.AlreadyExists, "Resource already exists")
	}

	data.SetID(primitive.NewObjectID())
	return data, nil
}

func (r *baseRepository[T]) FindOne(ctx context.Context, filters map[string]any) (T, error) {
	var zero T
	filter := MakeFilter(filters)
	result := r.Collection.FindOne(context.Background(), filter)

	if result.Err() != nil {
		return zero, result.Err()
	}

	var data T
	err := result.Decode(&data)
	if err != nil {
		return zero, err
	}

	return data, nil
}
func (r *baseRepository[T]) FindAll(ctx context.Context, filters map[string]any) ([]T, error) {
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

func (r *baseRepository[T]) Update(ctx context.Context, id primitive.ObjectID, data T) (T, error) {
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

func (r *baseRepository[T]) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := MakeFilter(map[string]any{"_id": id})
	_, err := r.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *baseRepository[T]) Count(ctx context.Context, filters map[string]any) (int64, error) {
	filter := MakeFilter(filters)
	count, err := r.Collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *baseRepository[T]) Transaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error) {
	session, err := r.Collection.Database().Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	wrappedFn := func(sc mongo.SessionContext) (any, error) {
		data, err := fn(sc)
		return data, err
	}

	return session.WithTransaction(ctx, wrappedFn)
}
