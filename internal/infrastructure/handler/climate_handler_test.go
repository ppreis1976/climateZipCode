package handler

import (
	"climateZpCode/internal/business/model"
	mockZipCode "climateZpCode/internal/infrastructure/service/mock"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClimateHandlerIntegration(t *testing.T) {
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

	// Configurar o servidor de teste
	handler := func(w http.ResponseWriter, r *http.Request) {
		ClimateHandler(w, r)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	tests := []struct {
		name           string
		method         string
		zipCode        string
		mockSetup      func() []*mock.Call
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Method not allowed",
			method: http.MethodPost,
			mockSetup: func() []*mock.Call {
				return nil
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed",
		},
		{
			name:   "Missing zip_code header",
			method: http.MethodGet,
			mockSetup: func() []*mock.Call {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "zip_code header is missing\n",
		},
		{
			name:    "Invalid zip code",
			method:  http.MethodGet,
			zipCode: "invalid",
			mockSetup: func() []*mock.Call {
				z := mockZipCodeService.On("Get", mock.Anything, model.ZipCodeID("invalid")).Return(model.ZipCode{}, errors.New(model.ErrZipCodeIDInvalid))
				return []*mock.Call{z}
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   model.ErrZipCodeIDInvalid + "\n",
		},
		{
			name:    "Zip code not found",
			method:  http.MethodGet,
			zipCode: "78890000",
			mockSetup: func() []*mock.Call {
				z := mockZipCodeService.On("Get", mock.Anything, model.ZipCodeID("03276010")).Return(model.ZipCode{}, nil)
				return []*mock.Call{z}
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   model.ErrZipCodeNotFound + "\n",
		},
		{
			name:    "Successful request",
			method:  http.MethodGet,
			zipCode: "03276010",
			mockSetup: func() []*mock.Call {
				z := mockZipCodeService.On("Get", mock.Anything, model.ZipCodeID("03276010")).Return(&model.ZipCode{Localidade: "local"}, nil)
				w := mockWeatherService.On("Get", mock.Anything, model.WeatherID("local")).Return(&weather, nil)
				return []*mock.Call{z, w}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"temp_C":16.4, "temp_F":61.519999999999996, "temp_K":289.4}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.mockSetup()

			req, err := http.NewRequest(tt.method, server.URL+"/climate", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.zipCode != "" {
				req.Header.Set(ZipCode, tt.zipCode)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				var body map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}

				expectedBody := make(map[string]interface{})
				if err := json.Unmarshal([]byte(tt.expectedBody), &expectedBody); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, expectedBody, body)
			}

			for _, call := range a {
				call.Unset()
			}

			mockZipCodeService.AssertExpectations(t)
			mockWeatherService.AssertExpectations(t)
		})
	}
}
