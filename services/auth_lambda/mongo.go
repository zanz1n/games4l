package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() error {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DATABASE_NAME")

	bcryptSaltLen, err := strconv.Atoi(os.Getenv("BCRYPT_SALT_LEN"))

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mongoConnStartTime := time.Now()

	serverApi := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return err
	}

	dba = user.NewUserService(client, &user.Config{
		MongoDbName:      dbName,
		BcryptSaltLength: bcryptSaltLen,
		JwtExpiryTime:    21 * 24 * time.Hour, // 21 days
	})

	logger.Info("Connected to mongodb, handshake took %v", time.Since(mongoConnStartTime))

	return nil
}
