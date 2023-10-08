package dashboard

import (
	"context"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/internal/repository"
)

type service struct {
	appRepository            repository.App
	shipRepository           repository.Ship
	pairingRequestRepository repository.PairingRequest
}

type Service interface {
	CountShip(ctx context.Context) (int64, error)
	GetShipsInBatch(ctx context.Context, start int, end int) ([]dto.ShipWebsocketResponse, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		appRepository:            f.AppRepository,
		shipRepository:           f.ShipRepository,
		pairingRequestRepository: f.PairingRequestRepository,
	}
}

func (s *service) CountShip(ctx context.Context) (int64, error) {
	res, err := s.shipRepository.CountShip(ctx)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (s *service) GetShipsInBatch(ctx context.Context, start int, end int) ([]dto.ShipWebsocketResponse, error) {
	res, is_update, err := s.shipRepository.ShipInBatch(ctx, start, end)
	if err != nil {
		return nil, err
	}

	var data []dto.ShipWebsocketResponse
	for _, e := range *res {
		data = append(data, dto.ShipWebsocketResponse{
			IsUpdate: is_update,
			ShipID:   e.ID,
			ShipName: e.Name,
			DeviceID: e.DeviceID,
			OnGround: e.OnGround,
			Lat:      e.CurrentLat,
			Long:     e.CurrentLong,
		})
	}

	return data, nil
}
