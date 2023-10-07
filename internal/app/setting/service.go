package setting

import (
	"context"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/internal/model"
	"simpel-api/internal/repository"
	"simpel-api/pkg/constants"
)

type service struct {
	AppSettingRepository  repository.AppSetting
	AppGeofenceRepository repository.AppGeofence
}

type Service interface {
	CreateOrUpdate(ctx context.Context, payload dto.PayloadStoreSetting) error
	GetSetting(ctx context.Context) (dto.GetDataSetting, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		AppSettingRepository:  f.AppSettingRepository,
		AppGeofenceRepository: f.AppGeofenceRepository,
	}
}

func (s *service) GetSetting(ctx context.Context) (dto.GetDataSetting, error) {
	appsetting, err := s.AppSettingRepository.FindLatest(ctx, "harbour_code, harbour_name, mode, apk_min_version, interval, range, apk_download_link")
	if err != nil {
		return dto.GetDataSetting{}, constants.NotFoundDataAppSetting
	}

	data := dto.GetDataSetting{
		HarbourCode:     appsetting.HarbourCode,
		HarbourName:     appsetting.HarbourName,
		Mode:            appsetting.Mode.String(),
		ApkMinVersion:   appsetting.ApkMinVersion,
		Interval:        appsetting.Interval,
		Range:           appsetting.Range,
		ApkDownloadLink: appsetting.ApkDownloadLink,
	}

	return data, nil
}

func (s *service) CreateOrUpdate(ctx context.Context, payload dto.PayloadStoreSetting) error {
	appsetting, err := s.AppSettingRepository.FindLatest(ctx, "id")

	if err != nil {
		dataStore := model.AppSetting{
			HarbourCode:     payload.HarbourCode,
			HarbourName:     payload.HarbourName,
			Mode:            model.ModeType(payload.Mode),
			ApkMinVersion:   payload.ApkMinVersion,
			Interval:        payload.Interval,
			Range:           payload.Range,
			ApkDownloadLink: payload.ApkDownloadLink,
		}

		s.AppSettingRepository.Store(ctx, dataStore)
	} else {
		update := model.AppSetting{
			HarbourCode:     payload.HarbourCode,
			HarbourName:     payload.HarbourName,
			ApkMinVersion:   payload.ApkMinVersion,
			Mode:            model.ModeType(payload.Mode),
			Interval:        payload.Interval,
			Range:           payload.Range,
			ApkDownloadLink: payload.ApkDownloadLink,
		}

		err = s.AppSettingRepository.Update(ctx, &update, "harbour_code,harbour_name,mode,apk_min_version,interval,range,apk_download_link,updated_at", "harbour_code = ?", appsetting.HarbourCode)
		if err != nil {
			return constants.ErrorUpdateAppSetting
		}
	}

	if payload.Geofence != nil {
		s.AppGeofenceRepository.DeleteAll(ctx)
		for _, geofence := range payload.Geofence {
			store := model.AppGeofence{
				Long: geofence.Lat,
				Lat:  geofence.Lat,
			}

			s.AppGeofenceRepository.Store(ctx, store)
		}
	}

	return nil
}
