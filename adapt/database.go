package adapt

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"log"
)

type Database struct {
	MongoClient  *mongo.Client
	DatabaseName string
	Log          *zap.Logger
}

func (app *application) buildMongoClient(ctx context.Context) *mongo.Client {
	clientOpts := app.config.GetClientOptions()
	mongoClient, err := mongo.NewClient(clientOpts)
	err = mongoClient.Connect(ctx)
	if err != nil {
		app.logger.Info("msg", zap.String("msg", "failed to connect to database"))
		log.Fatal(err)
	}
	return mongoClient
}
