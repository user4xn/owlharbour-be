package ship

import (
	"context"
	"fmt"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/internal/model"
	"simpel-api/internal/repository"
	"simpel-api/pkg/helper"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type service struct {
	appRepository            repository.App
	shipRepository           repository.Ship
	pairingRequestRepository repository.PairingRequest
}

type Service interface {
	PairingShip(ctx context.Context, request dto.PairingRequest) error
	PairingRequestList(ctx context.Context, request dto.PairingListParam) ([]dto.PairingRequestResponse, error)
	PairingAction(ctx context.Context, request dto.PairingActionRequest) error
	ShipByDevice(ctx context.Context, DeviceID string) (*dto.ShipDetailResponse, error)
	ShipList(ctx context.Context, request dto.ShipListParam) ([]dto.ShipResponse, error)
	RecordLocationShip(ctx context.Context, request dto.ShipRecordRequest) error
}

func NewService(f *factory.Factory) Service {
	return &service{
		appRepository:            f.AppRepository,
		shipRepository:           f.ShipRepository,
		pairingRequestRepository: f.PairingRequestRepository,
	}
}

func (s *service) PairingShip(ctx context.Context, request dto.PairingRequest) error {
	appInfo, err := s.appRepository.AppInfo(ctx)
	if err != nil {
		return err
	}

	if request.HarbourCode != appInfo.HarbourCode {
		return fmt.Errorf("unable to pair a device, invalid harbour code")
	}

	err = s.pairingRequestRepository.StorePairingRequests(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) PairingRequestList(ctx context.Context, request dto.PairingListParam) ([]dto.PairingRequestResponse, error) {
	res, err := s.pairingRequestRepository.PairingRequestList(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) PairingAction(ctx context.Context, request dto.PairingActionRequest) error {

	res, err := s.pairingRequestRepository.UpdatedPairingStatus(ctx, request)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid pairing_id")
		}
		return err
	}

	pairingData := dto.PairingRequestResponse{
		ID:              res.ID,
		ShipName:        res.ShipName,
		Phone:           res.Phone,
		ResponsibleName: res.ResponsibleName,
		DeviceID:        res.DeviceID,
		FirebaseToken:   res.FirebaseToken,
		Status:          res.Status,
		CreatedAt:       res.CreatedAt,
	}

	appInfo, err := s.appRepository.AppInfo(ctx)
	if err != nil {
		return err
	}

	if request.Status == "approved" {
		err = s.shipRepository.StoreNewShip(ctx, pairingData)
		if err != nil {
			return err
		}

		notificationData := map[string]interface{}{
			"title": "SIMPEL - PAIRING APPROVED",
			"body":  "Pairing request anda telah disetujui, kini device kapal anda sudah terhubung dengan Pelabuhan " + appInfo.HarbourName,
		}
		tokens := []string{res.FirebaseToken}

		_, err := helper.PushNotification(notificationData, tokens)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		notificationData := map[string]interface{}{
			"title": "SIMPEL - PAIRING REJECTED",
			"body":  "Mohon maaf pairing request device anda dengan Pelabuhan " + appInfo.HarbourName + "ditolak, anda dapat mengajukan kembali di lain waktu",
		}
		tokens := []string{res.FirebaseToken}

		_, err := helper.PushNotification(notificationData, tokens)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (s *service) ShipList(ctx context.Context, request dto.ShipListParam) ([]dto.ShipResponse, error) {
	res, err := s.shipRepository.ShipList(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) ShipByDevice(ctx context.Context, DeviceID string) (*dto.ShipDetailResponse, error) {
	res, err := s.shipRepository.ShipByDevice(ctx, DeviceID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) RecordLocationShip(ctx context.Context, request dto.ShipRecordRequest) error {
	ship, err := s.shipRepository.ShipByDevice(ctx, request.DeviceID)
	if err != nil {
		return err
	}

	appInfo, err := s.appRepository.AppInfo(ctx)
	if err != nil {
		return err
	}

	lat, err := strconv.ParseFloat(request.Lat, 64)
	if err != nil {
		return err
	}

	long, err := strconv.ParseFloat(request.Long, 64)
	if err != nil {
		return err
	}

	coord := [2]float64{lat, long}

	polygonData, err := s.appRepository.GetPolygon(ctx)
	if err != nil {
		return err
	}

	var polygon [][]float64
	for _, geo := range polygonData {
		lat, err := strconv.ParseFloat(geo.Lat, 64)
		if err != nil {
			return err
		}
		long, err := strconv.ParseFloat(geo.Long, 64)
		if err != nil {
			return err
		}
		polygon = append(polygon, []float64{lat, long})
	}

	polygon2D := convertPolygon(polygon)

	isInside := helper.StatusCheck(coord, polygon2D)

	var isWater bool
	var status string
	currentTime := time.Now()
	formattedTimeNotification := currentTime.Format("060102-1504")

	if isInside {
		lastLogs, _ := s.shipRepository.GetLastDockedLog(ctx, ship.ID)

		if lastLogs == nil || (lastLogs != nil && lastLogs.Status != "checkin") {
			dockedLog := dto.ShipDockedLogStore{
				ShipID: ship.ID,
				Lat:    request.Lat,
				Long:   request.Long,
				Status: "checkin",
			}

			if err := s.shipRepository.StoreDockedLog(ctx, dockedLog); err != nil {
				return err
			}

			notificationData := map[string]interface{}{
				"title": "SIMPEL - CHECK IN SUCCESS",
				"body":  "Berhasil CHECK-IN Pelabuhan " + appInfo.HarbourName + " " + formattedTimeNotification,
			}
			tokens := []string{ship.FirebaseToken}

			_, err := helper.PushNotification(notificationData, tokens)
			if err != nil {
				fmt.Println(err)
			}

			isWater = true
			status = "checkin"
		} else {
			isWater, err = helper.IsWater(lat, long)
			if err != nil {
				return err
			}

			status = "checkin"
		}
	} else {
		lastLogs, err := s.shipRepository.GetLastDockedLog(ctx, ship.ID)
		if err != nil {
			return err
		}

		if lastLogs != nil && lastLogs.Status == "checkin" {
			if ship.OnGround != 1 {
				isWater, err = helper.IsWater(lat, long)
				if err != nil {
					return err
				}

				if isWater {
					dockedLog := dto.ShipDockedLogStore{
						ShipID: ship.ID,
						Lat:    request.Lat,
						Long:   request.Long,
						Status: "checkout",
					}

					if err := s.shipRepository.StoreDockedLog(ctx, dockedLog); err != nil {
						return err
					}

					notificationData := map[string]interface{}{
						"title": "SIMPEL - CHECK OUT SUCCESS",
						"body":  "Berhasil CHECK-OUT Pelabuhan " + appInfo.HarbourName + " " + formattedTimeNotification,
					}
					tokens := []string{ship.FirebaseToken}

					_, err := helper.PushNotification(notificationData, tokens)
					if err != nil {
						fmt.Println(err)
					}
					status = "checkout"
				} else {
					status = "out of scope"
				}
			} else {
				isWater = false
				status = ship.Status
			}
		} else {
			isWater, err = helper.IsWater(lat, long)
			if err != nil {
				return err
			}

			status = "out of scope"
		}
	}
	fmt.Println(isWater, "is water")
	sll := dto.ShipLocationLogStore{
		ShipID:   ship.ID,
		Lat:      request.Lat,
		Long:     request.Long,
		IsMocked: request.IsMocked,
		OnGround: func() int {
			if isWater {
				return 0
			}
			return 1
		}(),
	}

	if err := s.shipRepository.StoreLocationLog(ctx, sll); err != nil {
		return err
	}

	setID := model.Common{
		ID: ship.ID,
	}

	shipUpdate := model.Ship{
		Common:      setID,
		Status:      model.ShipStatus(status),
		CurrentLat:  request.Lat,
		CurrentLong: request.Long,
		OnGround: func() int {
			if isWater {
				return 0
			}
			return 1
		}(),
	}

	if err := s.shipRepository.UpdateShip(ctx, shipUpdate); err != nil {
		return err
	}

	return nil
}

func convertPolygon(polygon [][]float64) [][2]float64 {
	result := make([][2]float64, len(polygon))
	for i, coord := range polygon {
		result[i] = [2]float64{coord[0], coord[1]}
	}
	return result
}
