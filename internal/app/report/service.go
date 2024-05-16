package report

import (
	"context"
	"owlharbour-api/internal/dto"
	"owlharbour-api/internal/factory"
	"owlharbour-api/internal/repository"
)

type service struct {
	shipRepository repository.Ship
}

type Service interface {
	ShipDocking(ctx context.Context, request dto.ReportShipDockedParam) ([]dto.ReportShipDockingResponse, error)
	ShipFraud(ctx context.Context, request dto.ReportShipLocationParam) ([]dto.ReportShipLocationResponse, error)
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

func (s *service) ShipFraud(ctx context.Context, request dto.ReportShipLocationParam) ([]dto.ReportShipLocationResponse, error) {
	res, err := s.shipRepository.ReportShipFraud(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}
