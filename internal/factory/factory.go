package factory

type Factory struct {
	// SomeRepository repository.Event
}

func NewFactory() *Factory {
	// Check db connection
	// db := database.GetConnection()
	return &Factory{
		// Pass the db connection to the repository package for database query calling
		// SomeRepository: repository.NewSomeRepository(db),
		// Assign the appropriate implementation of the ReturInsightRepository
	}
}
