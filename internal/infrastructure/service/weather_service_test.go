package service

import (
	mockService "climateZpCode/internal/infrastructure/service/mock"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeatherService_Get(t *testing.T) {
	mockClient := new(mockService.MockHTTPClient)
	service := NewWeatherService(mockClient)

	t.Run("success", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"location":{"name":"London","region":"City of London, Greater London","country":"United Kingdom","lat":51.5171,"lon":-0.1062,"tz_id":"Europe/London","localtime_epoch":1728858087,"localtime":"2024-10-13 23:21"},"current":{"last_updated_epoch":1728857700,"last_updated":"2024-10-13 23:15","temp_c":8.4,"temp_f":47.1,"is_day":0,"condition":{"text":"Clear","icon":"//cdn.weatherapi.com/weather/64x64/night/113.png","code":1000},"wind_mph":3.6,"wind_kph":5.8,"wind_degree":86,"wind_dir":"E","pressure_mb":1018.0,"pressure_in":30.06,"precip_mm":0.03,"precip_in":0.0,"humidity":76,"cloud":0,"feelslike_c":7.7,"feelslike_f":45.9,"windchill_c":8.7,"windchill_f":47.6,"heatindex_c":9.2,"heatindex_f":48.6,"dewpoint_c":2.3,"dewpoint_f":36.1,"vis_km":10.0,"vis_miles":6.0,"uv":0.0,"gust_mph":5.1,"gust_kph":8.2}}`)),
		}
		mockClient.On("Get", "http://api.weatherapi.com/v1/current.json?key=6ce46037f7034601a4e231757240909&q=London&aqi=no").Return(mockResponse, nil)

		weather, err := service.Get(context.Background(), "London")
		fmt.Sprintln(weather)
		assert.NoError(t, err)
		assert.Equal(t, "London", weather.Location.Name)
		assert.Equal(t, 8.4, weather.Current.TempC)
		assert.Equal(t, 47.1, weather.Current.TempF)
	})

	t.Run("network error", func(t *testing.T) {
		mockClient.On("Get", "http://api.weatherapi.com/v1/current.json?key=6ce46037f7034601a4e231757240909&q=Unknown&aqi=no").Return(nil, errors.New("network error"))

		_, err := service.Get(context.Background(), "Unknown")
		assert.Error(t, err)
		assert.Equal(t, "network error", err.Error())
	})

	t.Run("non-OK status", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(strings.NewReader(``)),
		}
		mockClient.On("Get", "http://api.weatherapi.com/v1/current.json?key=6ce46037f7034601a4e231757240909&q=ErrorCity&aqi=no").Return(mockResponse, nil)

		_, err := service.Get(context.Background(), "ErrorCity")
		assert.Error(t, err)
		assert.Equal(t, "failed to fetch data from weatherapi", err.Error())
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`invalid json`)),
		}
		mockClient.On("Get", "http://api.weatherapi.com/v1/current.json?key=6ce46037f7034601a4e231757240909&q=InvalidJSONCity&aqi=no").Return(mockResponse, nil)

		_, err := service.Get(context.Background(), "InvalidJSONCity")
		assert.Error(t, err)
	})
}
