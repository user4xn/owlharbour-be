package factory

import (
	"owlharbour-api/database"
	"owlharbour-api/internal/rabbitmq"
	"owlharbour-api/internal/repository"
	"owlharbour-api/pkg/util"

	"github.com/redis/go-redis/v9"
)

type Factory struct {
	AppRepository            repository.App
	ShipRepository           repository.Ship
	PairingRequestRepository repository.PairingRequest
	UserRepository           repository.User
	RabbitMqRepository       repository.RabbitMq
}

func NewFactory() *Factory {
	db := database.GetConnection()
	conn, ch := rabbitmq.CreateConnection()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     util.GetEnv("REDIS_URL", "localhost") + ":" + util.GetEnv("REDIS_PORT", "6379"),
		Password: util.GetEnv("REDIS_PASS", ""),
		DB:       0,
	})

	return &Factory{
		// Pass the db connection to the repository package for database query calling
		AppRepository:            repository.NewAppRepository(db, redisClient),
		ShipRepository:           repository.NewShipRepository(db, redisClient),
		PairingRequestRepository: repository.NewPairingRequestRepository(db, redisClient),
		UserRepository:           repository.NewUserRepository(db, redisClient),
		RabbitMqRepository:       repository.NewRabbitMqRepository(conn, ch),
		// Assign the appropriate implementation of the ReturInsightRepository
	}
}
