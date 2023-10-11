package report

import (
	"context"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/internal/repository"
)

type service struct {
	shipRepository repository.Ship
}

type Service interface {
	ShipDocking(ctx context.Context, request dto.ReportShipDockedParam) ([]dto.ReportShipDockingResponse, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		shipRepository: f.ShipRepository,
	}
}

func (s *service) ShipDocking(ctx context.Context, request dto.ReportShipDockedParam) ([]dto.ReportShipDockingResponse, error) {
	res, err := s.shipRepository.ReportShipDocking(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}
