package database

import (
	"context"
	"sync"
	"time"

	"github.com/itsanindyak/go-jwt/config"
	"github.com/itsanindyak/go-jwt/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	clientOnce sync.Once
)

func DbConnect() *mongo.Client {
	clientOnce.Do(func() {

		MongoURL := config.MONGO_URI

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI(MongoURL).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

		c, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			logger.Fatal("MongoDB connection error: " + err.Error())
		}

		if err := c.Ping(ctx, nil); err != nil {
			logger.Fatal("❌ MongoDB ping error:" + err.Error())
		}

		logger.Success("✅ Connected to MongoDB")
		client = c
	})

	return client
}

func GetCollection(ctx context.Context, name string) *mongo.Collection {
	return DbConnect().Database("testdb").Collection(name)
}
