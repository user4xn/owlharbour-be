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
	GetStatistic(ctx context.Context) (*dto.DashboardStatisticResponse, error)
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
			Geo:      []string{e.CurrentLat, e.CurrentLong},
		})
	}

	return data, nil
}

func (s *service) GetStatistic(ctx context.Context) (*dto.DashboardStatisticResponse, error) {
	countShip, err := s.shipRepository.CountShip(ctx)
	if err != nil {
		return nil, err
	}

	countStatistic, err := s.shipRepository.CountStatistic(ctx)
	if err != nil {
		return nil, err
	}

	res := dto.DashboardStatisticResponse{
		TotalShip:     int(countShip),
		TotalCheckin:  int(countStatistic[0]),
		TotalCheckout: int(countStatistic[1]),
		TotalFraud:    int(countStatistic[2]),
	}

	return &res, nil
}
