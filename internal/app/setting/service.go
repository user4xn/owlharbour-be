package setting

import (
	"context"
	"owlharbour-api/internal/dto"
	"owlharbour-api/internal/factory"
	"owlharbour-api/internal/model"
	"owlharbour-api/internal/repository"
	"owlharbour-api/pkg/constants"
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
	appsetting, err := s.AppRepository.FindLatestSetting(ctx, "harbour_code, harbour_name, mode, interval, range, admin_contact")
	if err != nil {
		return dto.GetDataSetting{}, constants.NotFoundDataAppSetting
	}

	getGeofance, err := s.AppRepository.GetPolygon(ctx)
	if err != nil {
		data := dto.GetDataSetting{
			HarbourCode:  appsetting.HarbourCode,
			HarbourName:  appsetting.HarbourName,
			Mode:         appsetting.Mode.String(),
			Interval:     appsetting.Interval,
			Range:        appsetting.Range,
			AdminContact: appsetting.AdminContact,
			Geofences:    nil,
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

	data := dto.GetDataSetting{
		HarbourCode:  appsetting.HarbourCode,
		HarbourName:  appsetting.HarbourName,
		Mode:         appsetting.Mode.String(),
		Interval:     appsetting.Interval,
		Range:        appsetting.Range,
		AdminContact: appsetting.AdminContact,
		Geofences:    geofences,
	}

	return data, nil
}

func (s *service) GetSettingWeb(ctx context.Context) (dto.GetDataSettingWeb, error) {
	appsetting, err := s.AppRepository.FindLatestSetting(ctx, "harbour_code, harbour_name, mode, interval, range, admin_contact")
	if err != nil {
		return dto.GetDataSettingWeb{}, err
	}

	getGeofance, err := s.AppRepository.GetPolygon(ctx)
	if err != nil {
		data := dto.GetDataSettingWeb{
			HarbourCode:  appsetting.HarbourCode,
			HarbourName:  appsetting.HarbourName,
			Mode:         appsetting.Mode.String(),
			Interval:     appsetting.Interval,
			Range:        appsetting.Range,
			AdminContact: appsetting.AdminContact,
			Geofences:    nil,
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
		HarbourCode:  appsetting.HarbourCode,
		HarbourName:  appsetting.HarbourName,
		Mode:         appsetting.Mode.String(),
		Interval:     appsetting.Interval,
		Range:        appsetting.Range,
		AdminContact: appsetting.AdminContact,
		Geofences:    geofences,
	}

	return data, nil
}
func (s *service) CreateOrUpdate(ctx context.Context, payload dto.PayloadStoreSetting) error {
	appsetting, err := s.AppRepository.FindLatestSetting(ctx, "harbour_code")
	if err != nil {
		dataStore := model.AppSetting{
			HarbourCode:  payload.HarbourCode,
			HarbourName:  payload.HarbourName,
			Mode:         model.ModeType(payload.Mode),
			Interval:     payload.Interval,
			Range:        payload.Range,
			AdminContact: payload.AdminContact,
		}

		s.AppRepository.StoreSetting(ctx, dataStore)
	} else {
		update := model.AppSetting{
			HarbourCode:  payload.HarbourCode,
			HarbourName:  payload.HarbourName,
			Mode:         model.ModeType(payload.Mode),
			Interval:     payload.Interval,
			Range:        payload.Range,
			AdminContact: payload.AdminContact,
		}

		s.AppRepository.UpsertSetting(ctx, &update, "harbour_code,harbour_name,mode,interval,range,admin_contact,updated_at", "harbour_code = ?", appsetting.HarbourCode)
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
