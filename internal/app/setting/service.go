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
	AppRepository repository.App
}

type Service interface {
	CreateOrUpdate(ctx context.Context, payload dto.PayloadStoreSetting) error
	GetSetting(ctx context.Context) (dto.GetDataSetting, error)
	GetSettingWeb(ctx context.Context) (dto.GetDataSettingWeb, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		AppRepository: f.AppRepository,
	}
}

func (s *service) GetSetting(ctx context.Context) (dto.GetDataSetting, error) {
	appsetting, err := s.AppRepository.FindLatestSetting(ctx, "harbour_code, harbour_name, mode, apk_min_version, interval, range, apk_download_link")
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

func (s *service) GetSettingWeb(ctx context.Context) (dto.GetDataSettingWeb, error) {
	appsetting, err := s.AppRepository.FindLatestSetting(ctx, "harbour_code, harbour_name, mode, apk_min_version, interval, range, apk_download_link")
	if err != nil {
		return dto.GetDataSettingWeb{}, err
	}

	getGeofance, err := s.AppRepository.GetPolygon(ctx)
	if err != nil {
		data := dto.GetDataSettingWeb{
			HarbourCode:     appsetting.HarbourCode,
			HarbourName:     appsetting.HarbourName,
			Mode:            appsetting.Mode.String(),
			ApkMinVersion:   appsetting.ApkMinVersion,
			Interval:        appsetting.Interval,
			Range:           appsetting.Range,
			ApkDownloadLink: appsetting.ApkDownloadLink,
			Geofences:       nil,
		}
		return data, nil
	}

	geofences := []dto.AppGeofence{}
	for _, gf := range getGeofance {
		dataGeofance := dto.AppGeofence{
			Long: gf.Long,
			Lat:  gf.Lat,
		}
		geofences = append(geofences, dataGeofance)
	}
	data := dto.GetDataSettingWeb{
		HarbourCode:     appsetting.HarbourCode,
		HarbourName:     appsetting.HarbourName,
		Mode:            appsetting.Mode.String(),
		ApkMinVersion:   appsetting.ApkMinVersion,
		Interval:        appsetting.Interval,
		Range:           appsetting.Range,
		ApkDownloadLink: appsetting.ApkDownloadLink,
		Geofences:       geofences,
	}

	return data, nil
}
func (s *service) CreateOrUpdate(ctx context.Context, payload dto.PayloadStoreSetting) error {
	appsetting, err := s.AppRepository.FindLatestSetting(ctx, "harbour_code")
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

		s.AppRepository.StoreSetting(ctx, dataStore)
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

		s.AppRepository.UpsertSetting(ctx, &update, "harbour_code,harbour_name,mode,apk_min_version,interval,range,apk_download_link,updated_at", "harbour_code = ?", appsetting.HarbourCode)
	}

	if payload.Geofence != nil {
		s.AppRepository.DeleteAllGeofence(ctx)
		for _, geofence := range payload.Geofence {
			store := model.AppGeofence{
				Long: geofence.Long,
				Lat:  geofence.Lat,
			}

			s.AppRepository.StoreGeofence(ctx, store)
		}
	}

	return nil
}
