package usercase

import (
	"climateZpCode/internal/business/gateway"
	"climateZpCode/internal/business/model"
	"context"
)

type (
	ClimateUserCase interface {
		Get(zipCodeID model.ZipCodeID) (model.Climate, error)
	}

	climateUserCase struct {
		zipCodeService gateway.ZipCodeService
		weatherService gateway.WheaterService
	}
)

func NewClimateUserCase(zipCodeService gateway.ZipCodeService, weatherService gateway.WheaterService) ClimateUserCase {
	return climateUserCase{
		zipCodeService: zipCodeService,
		weatherService: weatherService,
	}
}

func (u climateUserCase) Get(zipCodeID model.ZipCodeID) (model.Climate, error) {
	if err := zipCodeID.Validate(); err != nil {
		return model.Climate{}, err
	}

	ctx := context.Background()
	z, err := u.zipCodeService.Get(ctx, zipCodeID)
	if err != nil {
		return model.Climate{}, err
	}

	weatherID := model.WeatherID(z.Localidade)
	weather, err := u.weatherService.Get(ctx, weatherID)
	if err != nil {
		return model.Climate{}, err
	}

	var climate model.Climate

	climate.TempK = weather.Current.TempC
	climate.TempC = weather.Current.TempC
	climate.TempF = weather.Current.TempF

	return climate, nil
}
