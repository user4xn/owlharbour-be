package ship

import (
	"context"
	"fmt"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/internal/repository"
	"simpel-api/pkg/helper"

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