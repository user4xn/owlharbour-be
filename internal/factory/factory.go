package factory

import (
	"simpel-api/database"
	"simpel-api/internal/repository"
)

type Factory struct {
	UserRepository repository.UserInterface
}

func NewFactory() *Factory {
	// Check db connection
	db := database.GetConnection()
	return &Factory{
		// Pass the db connection to the repository package for database query calling
		// SomeRepository: repository.NewSomeRepository(db),
		// Assign the appropriate implementation of the ReturInsightRepository
		UserRepository: repository.NewUserRepository(db),
	}
}
