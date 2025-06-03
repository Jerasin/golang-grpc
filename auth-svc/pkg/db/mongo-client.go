package db

import (
	"auth-svc/pkg/config"
	"context"
	"log"
	"time"

	"github.com/goforj/godump"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnection() (*mongo.Client, context.CancelFunc, error) {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed at config ", err)
	}

	godump.Dump(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		c.MongoUri).SetServerSelectionTimeout(5*time.
		Second))
	if err != nil {
		cancel()
		return nil, nil, err
	}
	// db := client.Database("books")
	return client, cancel, nil
}
