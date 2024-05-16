package ship

import (
	"context"
	"fmt"
	"math/rand"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/internal/model"
	"simpel-api/internal/repository"
	"simpel-api/pkg/helper"
	"simpel-api/pkg/log"
	"simpel-api/pkg/util"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type service struct {
	appRepository            repository.App
	shipRepository           repository.Ship
	pairingRequestRepository repository.PairingRequest
	userRepository           repository.User
	RabbitMqRepository       repository.RabbitMq
}

type Service interface {
	PairingShip(ctx context.Context, request dto.PairingRequest) error
	PairingRequestCount(ctx context.Context) (int64, error)
	PairingRequestList(ctx context.Context, request dto.PairingListParam) ([]dto.PairingRequestResponse, error)
	PairingAction(ctx context.Context, request dto.PairingActionRequest) error
	PairingDetailByUsername(ctx context.Context, username string) (*dto.DetailPairingResponse, error)
	ShipByAuth(ctx context.Context, authUser model.User) (*dto.ShipMobileDetailResponse, error)
	ShipList(ctx context.Context, request dto.ShipListParam) ([]dto.ShipResponse, error)
	RecordLocationShip(ctx context.Context, request dto.ShipRecordRequest) error
	UpdateShipDetail(ctx context.Context, request dto.ShipAddonDetailRequest) error
	ShipDetail(ctx context.Context, ShipID int) (*dto.ShipDetailResponse, error)
	ShipDockLog(ctx context.Context, request dto.ShipLogParam, shipOrDeviceID any) (*dto.ShipDockLogResponse, error)
	ShipLocationLog(ctx context.Context, request dto.ShipLogParam, shipOrDeviceID any) (*dto.ShipLocationLogResponse, error)
	RecordShipRabbit(ctx context.Context, request dto.ShipRecordRequest) error
}

func NewService(f *factory.Factory) Service {
	return &service{
		appRepository:            f.AppRepository,
		shipRepository:           f.ShipRepository,
		pairingRequestRepository: f.PairingRequestRepository,
		userRepository:           f.UserRepository,
		RabbitMqRepository:       f.RabbitMqRepository,
	}
}

func (s *service) RecordShipRabbit(ctx context.Context, request dto.ShipRecordRequest) error {
	publishRequest := dto.RabbitMqPublishRequest{
		Exchange:  util.GetEnv("RABBITMQ_EXCHANGE_SIMPEL_SHIP", ""),
		QueueName: "ShipRecordLog",
		Messages:  request,
	}
	go func() {
		err := s.RabbitMqRepository.Publish(ctx, publishRequest)
		if err != nil {
			fmt.Println("Failed to publish a message", zap.String("device id", request.DeviceID), zap.String("error", err.Error()))
		}
	}()

	return nil
}

func (s *service) PairingRequestCount(ctx context.Context) (int64, error) {
	countPairing, err := s.pairingRequestRepository.PairingRequestCount(ctx, "pending")
	if err != nil {
		return 0, err
	}

	return countPairing, nil
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
	idArray := strings.Split(request.PairingID, ",")
	wg := sync.WaitGroup{}
	var failures sync.Mutex

	for _, id := range idArray {
		wg.Add(1)

		idInt, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("invalid pairing_id: %v", err)
		}

		go func(idInt int, id string) error {
			defer wg.Done()
			defer failures.Unlock()
			failures.Lock()

			res, err := s.pairingRequestRepository.UpdatedPairingStatus(ctx, idInt, request.Status)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					err = fmt.Errorf("invalid pairing_id: %s", id)
				}

				log.Logging("Failed to update pairing status (ID:%s), Err: %s", id, err.Error()).Error()
				return err
			}

			appInfo, err := s.appRepository.AppInfo(ctx)
			if err != nil {
				return err
			}

			if request.Status == "approved" {
				currentTime := time.Now()

				dataStore := model.User{
					Name:            res.ShipName,
					Email:           "",
					Username:        res.Username,
					EmailVerifiedAt: &currentTime,
					Password:        res.Password,
					Role:            "User",
				}
				err := s.userRepository.Store(ctx, dataStore)
				if err != nil {
					log.Logging("Failed store ship user, Pairing ID: %d, Err: %s", idInt, err.Error()).Error()

					return err
				}

				user, err := s.userRepository.FindOne(ctx, "id", "username = ?", res.Username)
				if err != nil {
					log.Logging("Failed get ship user, Pairing ID: %d, Err: %s", idInt, err.Error()).Error()

					return err
				}

				pairingToShip := dto.PairingToNewShip{
					ShipName:        res.ShipName,
					ResponsibleName: res.ResponsibleName,
					DeviceID:        res.DeviceID,
					FirebaseToken:   res.FirebaseToken,
					UserID:          user.ID,
					Phone:           res.Phone,
				}

				err = s.shipRepository.StoreNewShip(ctx, pairingToShip)
				if err != nil {
					log.Logging("Failed store new ship, ID: %d, Err: %s", idInt, err.Error()).Error()

					return err
				}

				notificationData := map[string]interface{}{
					"title": "OWLHARBOUR - PAIRING APPROVED",
					"body":  "Your ship pairing registration was approved, now your device connected to " + appInfo.HarbourName + " Harbour",
				}
				tokens := []string{res.FirebaseToken}

				_, err = helper.PushNotification(notificationData, tokens)
				if err != nil {
					log.Logging("Failed send notification, Err: %s", err.Error()).Error()
				}
			} else {
				notificationData := map[string]interface{}{
					"title": "OWLHARBOUR - PAIRING REJECTED",
					"body":  "We really sorry, your ship pairing registration was rejected by " + appInfo.HarbourName + " Harbour, please try again later",
				}
				tokens := []string{res.FirebaseToken}

				_, err := helper.PushNotification(notificationData, tokens)
				if err != nil {
					log.Logging("Failed send notification, Err: %s", err.Error()).Error()
				}
			}

			return nil
		}(idInt, id)
	}

	wg.Wait()
	failures.Lock()
	defer failures.Unlock()
	return nil
}

func (s *service) ShipList(ctx context.Context, request dto.ShipListParam) ([]dto.ShipResponse, error) {
	res, err := s.shipRepository.ShipList(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) ShipByAuth(ctx context.Context, authUser model.User) (*dto.ShipMobileDetailResponse, error) {
	appInfo, err := s.appRepository.AppInfo(ctx)
	if err != nil {
		return nil, err
	}

	ship, err := s.shipRepository.ShipByAuth(ctx, authUser)
	if err != nil {
		return nil, err
	}

	res := &dto.ShipMobileDetailResponse{
		ID:              ship.ID,
		ShipName:        ship.ShipName,
		ResponsibleName: ship.ResponsibleName,
		DeviceID:        ship.DeviceID,
		CurrentLong:     ship.CurrentLong,
		CurrentLat:      ship.CurrentLat,
		FirebaseToken:   ship.FirebaseToken,
		Status:          string(ship.Status),
		OnGround:        ship.OnGround,
		CreatedAt:       ship.CreatedAt,
		HitMode:         appInfo.Mode,
		Range:           appInfo.Range,
		Interval:        appInfo.Interval,
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
				"title": "OWLHARBOUR - CHECK IN SUCCESS",
				"body":  "Ship was checkin-in into " + appInfo.HarbourName + " Harbour at " + formattedTimeNotification,
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
		lastLogs, _ := s.shipRepository.GetLastDockedLog(ctx, ship.ID)

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
						"title": "OWLHARBOUR - CHECK OUT SUCCESS",
						"body":  "Ship was checkin-out from " + appInfo.HarbourName + " Harbour at " + formattedTimeNotification,
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

func (s *service) UpdateShipDetail(ctx context.Context, request dto.ShipAddonDetailRequest) error {
	err := s.shipRepository.UpdateShipDetail(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ShipDetail(ctx context.Context, ShipID int) (*dto.ShipDetailResponse, error) {
	ship, err := s.shipRepository.ShipByID(ctx, ShipID)
	if err != nil {
		return nil, err
	}

	addonDetail, err := s.shipRepository.ShipAddonDetail(ctx, ShipID)
	if err != nil {
		log.Logging("Error Fethcing Addon Ship %s", err).Info()
	}

	randomDeg := -180 + rand.Float64()*(360)
	res := &dto.ShipDetailResponse{
		ID:              ship.ID,
		ShipName:        ship.Name,
		ResponsibleName: ship.ResponsibleName,
		DeviceID:        ship.DeviceID,
		Phone:           ship.Phone,
		DetailShip:      addonDetail,
		CurrentLong:     ship.CurrentLong,
		CurrentLat:      ship.CurrentLat,
		DegNorth:        randomDeg,
		FirebaseToken:   ship.FirebaseToken,
		Status:          string(ship.Status),
		OnGround:        ship.OnGround,
		CreatedAt:       ship.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return res, nil
}

func (s *service) PairingDetailByUsername(ctx context.Context, username string) (*dto.DetailPairingResponse, error) {
	res, err := s.pairingRequestRepository.PairingDetailByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) ShipDockLog(ctx context.Context, request dto.ShipLogParam, shipOrDeviceID any) (*dto.ShipDockLogResponse, error) {
	var id int

	switch shipOrDeviceID.(type) {
	case int:
		id = shipOrDeviceID.(int)
	case string:
		ship, err := s.shipRepository.ShipByDevice(ctx, shipOrDeviceID.(string))
		if err != nil {
			return nil, err
		}

		id = ship.ID
	}

	dockLogs, err := s.shipRepository.ShipDockedLogs(ctx, id, &request)
	if err != nil {
		return nil, err
	}

	res := &dto.ShipDockLogResponse{
		ID:          id,
		DockingLogs: dockLogs,
	}

	return res, nil
}

func (s *service) ShipLocationLog(ctx context.Context, request dto.ShipLogParam, shipOrDeviceID any) (*dto.ShipLocationLogResponse, error) {
	var id int

	switch shipOrDeviceID.(type) {
	case int:
		id = shipOrDeviceID.(int)
	case string:
		ship, err := s.shipRepository.ShipByDevice(ctx, shipOrDeviceID.(string))
		if err != nil {
			return nil, err
		}

		id = ship.ID
	}

	locationLogs, err := s.shipRepository.ShipLocationLogs(ctx, id, &request)
	if err != nil {
		return nil, err
	}

	res := &dto.ShipLocationLogResponse{
		ID:           id,
		LocationLogs: locationLogs,
	}

	return res, nil
}
