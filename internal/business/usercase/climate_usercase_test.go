package usercase

import (
	"climateZpCode/internal/business/model"
	mockZipCode "climateZpCode/internal/infrastructure/service/mock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_climateUserCase_Get(t *testing.T) {
	mockZipCodeService := new(mockZipCode.MockZipCodeService)
	mockWeatherService := new(mockZipCode.MockWheaterService)

	weather := model.Weather{
		Location: model.Location{
			Name:           "Sao Paulo",
			Region:         "Sao Paulo",
			Country:        "Brazil",
			Lat:            -23.5333,
			Lon:            -46.6167,
			TzId:           "America/Sao_Paulo",
			LocaltimeEpoch: 1728861669,
			Localtime:      "2024-10-13 20:21",
		},
		Current: model.Current{
			LastUpdatedEpoch: 1728861300,
			LastUpdated:      "2024-10-13 20:15",
			TempC:            17.3,
			TempF:            63.1,
			IsDay:            0,
			Condition: model.Condition{
				Text: "Partly cloudy",
				Icon: "//cdn.weatherapi.com/weather/64x64/night/116.png",
				Code: 1003,
			},
			WindMph:    11.9,
			WindKph:    19.1,
			WindDegree: 131,
			WindDir:    "SE",
			PressureMb: 1027.0,
			PressureIn: 30.33,
			PrecipMm:   0.02,
			PrecipIn:   0.0,
			Humidity:   88,
			Cloud:      75,
			FeelslikeC: 17.3,
			FeelslikeF: 63.1,
			WindchillC: 17.4,
			WindchillF: 63.3,
			HeatindexC: 17.4,
			HeatindexF: 63.3,
			DewpointC:  15.8,
			DewpointF:  60.5,
			VisKm:      10.0,
			VisMiles:   6.0,
			Uv:         0.0,
			GustMph:    13.6,
			GustKph:    22.0,
		},
	}

	tests := []struct {
		name           string
		zipCodeID      model.ZipCodeID
		mockSetup      func() []*mock.Call
		expectedResult model.Climate
		expectError    bool
	}{
		{
			name:      "Invalid zip code",
			zipCodeID: "invalid",
			mockSetup: func() []*mock.Call {
				return nil
			},
			expectedResult: model.Climate{},
			expectError:    true,
		},
		{
			name:      "Zip code not found",
			zipCodeID: "03276010",
			mockSetup: func() []*mock.Call {
				z := mockZipCodeService.On("Get", mock.Anything, model.ZipCodeID("03276010")).Return(&model.ZipCode{}, errors.New(model.ErrZipCodeNotFound)).Maybe()
				return []*mock.Call{z}
			},
			expectedResult: model.Climate{},
			expectError:    true,
		},
		{
			name:      "Weather service error",
			zipCodeID: "03276010",
			mockSetup: func() []*mock.Call {
				z := mockZipCodeService.On("Get", mock.Anything, model.ZipCodeID("03276010")).Return(&model.ZipCode{Localidade: "local"}, nil).Maybe()
				w := mockWeatherService.On("Get", mock.Anything, model.WeatherID("local")).Return(&model.Weather{}, errors.New("weather service error")).Maybe()
				return []*mock.Call{z, w}
			},
			expectedResult: model.Climate{},
			expectError:    true,
		},
		{
			name:      "Successful request",
			zipCodeID: "03276010",
			mockSetup: func() []*mock.Call {
				z := mockZipCodeService.On("Get", mock.Anything, model.ZipCodeID("03276010")).Return(&model.ZipCode{Localidade: "local"}, nil).Maybe()
				w := mockWeatherService.On("Get", mock.Anything, model.WeatherID("local")).Return(&weather, nil).Maybe()
				return []*mock.Call{z, w}
			},
			expectedResult: model.Climate{TempC: 17.3, TempF: 63.14, TempK: 290.3},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			u := climateUserCase{
				zipCodeService: mockZipCodeService,
				weatherService: mockWeatherService,
			}
			a := tt.mockSetup()

			result, err := u.Get(tt.zipCodeID)

			for _, call := range a {
				call.Unset()
			}

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)

			mockZipCodeService.AssertExpectations(t)
			mockWeatherService.AssertExpectations(t)
		})
	}
}
