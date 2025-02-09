package usercase

import (
	"climateZpCode/internal/business/model"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockZipCodeService é um mock do interface ZipCodeService
type MockZipCodeService struct {
	mock.Mock
}

func (m *MockZipCodeService) Get(ctx context.Context, zipCodeID model.ZipCodeID) (model.ZipCode, error) {
	args := m.Called(ctx, zipCodeID)
	return args.Get(0).(model.ZipCode), args.Error(1)
}

// MockWeatherService é um mock do interface WeatherService
type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) Get(ctx context.Context, weatherID model.WeatherID) (model.Weather, error) {
	args := m.Called(ctx, weatherID)
	return args.Get(0).(model.Weather), args.Error(1)
}
