package src

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() error {
	if dba != nil {
		return nil
	}

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DATABASE_NAME")

	bcryptSaltLen, err := strconv.Atoi(os.Getenv("BCRYPT_SALT_LEN"))

	if err != nil {
		return errors.New("failed to parse BCRYPT_SALT_LEN environment variable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mongoConnStartTime := time.Now()

	serverApi := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return errors.New("failed to connect to mongodb: " + err.Error())
	}

	dba = user.NewUserService(client, ap, &user.Config{
		MongoDbName:      dbName,
		BcryptSaltLength: bcryptSaltLen,
		JwtExpiryTime:    time.Hour,
	})

	logger.Info("Connected to mongodb, handshake took %v", time.Since(mongoConnStartTime))

	return nil
}
