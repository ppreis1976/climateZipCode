package service

import (
	mockService "climateZpCode/internal/infrastructure/service/mock"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClimateService_Get(t *testing.T) {
	mockClient := new(mockService.MockHTTPClient)
	service := NewClimateService(mockClient)

	t.Run("success", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"cep": "01001-000", "logradouro": "Praça da Sé"}`)),
		}
		mockClient.On("Get", "http://viacep.com.br/ws/01001000/json/").Return(mockResponse, nil)

		zipCode, err := service.Get(context.Background(), "01001000")
		assert.NoError(t, err)
		assert.Equal(t, "01001-000", string(zipCode.Cep))
		assert.Equal(t, "Praça da Sé", zipCode.Logradouro)
	})

	t.Run("network error", func(t *testing.T) {
		mockClient.On("Get", "http://viacep.com.br/ws/00000000/json/").Return(nil, errors.New("network error"))

		_, err := service.Get(context.Background(), "00000000")
		assert.Error(t, err)
		assert.Equal(t, "network error", err.Error())
	})

	t.Run("non-OK status", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(strings.NewReader(``)),
		}
		mockClient.On("Get", "http://viacep.com.br/ws/11111111/json/").Return(mockResponse, nil)

		_, err := service.Get(context.Background(), "11111111")
		assert.Error(t, err)
		assert.Equal(t, "failed to fetch data from ViaCEP", err.Error())
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`invalid json`)),
		}
		mockClient.On("Get", "http://viacep.com.br/ws/22222222/json/").Return(mockResponse, nil)

		_, err := service.Get(context.Background(), "22222222")
		assert.Error(t, err)
	})
}
