package dashboard

import (
	"context"
	"fmt"
	"math/rand"
	"owlharbour-api/internal/dto"
	"owlharbour-api/internal/factory"
	"owlharbour-api/internal/repository"
	"time"
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
	TerrainChart(ctx context.Context) (*dto.ShipTerrainResponse, error)
	LogsChart(ctx context.Context, startDate string, endDate string) (*dto.LogsStatisticResponse, error)
	LastestDockedShip(ctx context.Context, limit int) ([]dto.DashboardLastDockedShipResponse, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		appRepository:            f.AppRepository,
		shipRepository:           f.ShipRepository,
		pairingRequestRepository: f.PairingRequestRepository,
	}
}

func (s *service) LastestDockedShip(ctx context.Context, limit int) ([]dto.DashboardLastDockedShipResponse, error) {
	res, err := s.shipRepository.LastestDockedShip(ctx, limit)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) LogsChart(ctx context.Context, startDate string, endDate string) (*dto.LogsStatisticResponse, error) {
	checkin, err := s.shipRepository.CountShipByStatus(ctx, startDate, endDate, "checkin")
	if err != nil {
		return nil, err
	}

	checkout, err := s.shipRepository.CountShipByStatus(ctx, startDate, endDate, "checkout")
	if err != nil {
		return nil, err
	}

	outOfScope, err := s.shipRepository.CountShipByStatus(ctx, startDate, endDate, "out of scope")
	if err != nil {
		return nil, err
	}

	fraud, err := s.shipRepository.CountShipFraud(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	res := dto.LogsStatisticResponse{
		OutOfScope: outOfScope,
		CheckIN:    checkin,
		CheckOUT:   checkout,
		Fraud:      fraud,
	}

	return &res, nil
}

func (s *service) TerrainChart(ctx context.Context) (*dto.ShipTerrainResponse, error) {
	onGorund, err := s.shipRepository.CountShipByTerrain(ctx, 1)
	if err != nil {
		return nil, err
	}

	onWater, err := s.shipRepository.CountShipByTerrain(ctx, 0)
	if err != nil {
		return nil, err
	}

	res := dto.ShipTerrainResponse{
		OnGround: onGorund,
		OnWater:  onWater,
	}

	return &res, nil
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
		fmt.Println("Res was ? and ?", e.CurrentLat, e.CurrentLong)
		rand.Seed(time.Now().UnixNano())

		// Generate a random number between -180 and 180
		randomDeg := -180 + rand.Float64()*(360)
		data = append(data, dto.ShipWebsocketResponse{
			IsUpdate: is_update,
			ShipID:   e.ID,
			ShipName: e.Name,
			DeviceID: e.DeviceID,
			OnGround: e.OnGround,
			Geo:      []string{e.CurrentLong, e.CurrentLat},
			DegNorth: randomDeg,
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
