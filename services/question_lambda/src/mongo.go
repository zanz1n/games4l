package src

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/question"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() error {
	if dba != nil {
		return nil
	}

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DATABASE_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mongoConnStartTime := time.Now()

	serverApi := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return errors.New("Failed to connect to mongodb: " + err.Error())
	}

	dba = question.NewQuestionService(client, &question.Config{
		MongoDbName: dbName,
	})

	logger.Info("Connected to mongodb, handshake took %v", time.Since(mongoConnStartTime))

	return nil
}
