package inspection

import (
	"context"
	"fmt"
	"owlharbour-api/internal/dto"
	"owlharbour-api/internal/factory"
	"owlharbour-api/internal/repository"
)

type service struct {
	shipRepository repository.Ship
}

type Service interface {
	UpdateShipCheckup(ctx context.Context, request dto.ShipCheckupRequest, id int) error
	NeedCheckupShip(ctx context.Context, request dto.NeedCheckupShipParam) ([]dto.NeedCheckupShipResponse, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		shipRepository: f.ShipRepository,
	}
}

func (s *service) NeedCheckupShip(ctx context.Context, request dto.NeedCheckupShipParam) ([]dto.NeedCheckupShipResponse, error) {
	res, err := s.shipRepository.NeedCheckupShip(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) UpdateShipCheckup(ctx context.Context, request dto.ShipCheckupRequest, id int) error {
	log, err := s.shipRepository.FindOneDockedLog(ctx, "is_inspected, is_reported", "id = ?", id)
	if err != nil {
		return err
	}

	fmt.Println(request)

	if log.IsInspected != 0 && request.IsInspected {
		return fmt.Errorf("ship inspected already updated")
	}

	if log.IsReported != 0 && request.IsReported {
		return fmt.Errorf("ship reported already updated")
	}

	err = s.shipRepository.UpdateShipCheckup(ctx, request, id, log)
	if err != nil {
		return err
	}

	return nil
}
