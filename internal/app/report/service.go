package report

import (
	"simpel-api/internal/factory"
)

type service struct {
	// someRepository repository.RepositoryName
}

type Service interface {
	// FucntionName(req) (res)
}

func NewService(f *factory.Factory) Service {
	return &service{
		// someRepository: f.SomeRepository,
	}
}
