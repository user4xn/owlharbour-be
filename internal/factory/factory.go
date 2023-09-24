package factory

import (
	"simpel-api/database"
	"simpel-api/internal/repository"
)

type Factory struct {
	AppRepository            repository.App
	ShipRepository           repository.Ship
	PairingRequestRepository repository.PairingRequest
	UserRepository           repository.UserInterface
}

func NewFactory() *Factory {
	// Check db connection
	db := database.GetConnection()
	return &Factory{
		// Pass the db connection to the repository package for database query calling
		AppRepository:            repository.NewAppRepository(db),
		ShipRepository:           repository.NewShipRepository(db),
		PairingRequestRepository: repository.NewPairingRequestRepository(db),
		UserRepository:           repository.NewUserRepository(db),
		// Assign the appropriate implementation of the ReturInsightRepository
	}
}
