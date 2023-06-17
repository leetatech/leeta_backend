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

func (app *Application) buildMongoClient(ctx context.Context) *mongo.Client {
	clientOpts := app.Config.GetClientOptions()
	mongoClient, err := mongo.NewClient(clientOpts)
	err = mongoClient.Connect(ctx)
	if err != nil {
		app.Logger.Info("msg", zap.String("msg", "failed to connect to database"))
		log.Fatal(err)
	}
	return mongoClient
}
