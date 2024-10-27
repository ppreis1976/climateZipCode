package service

import (
	"climateZpCode/internal/business/gateway"
	"climateZpCode/internal/business/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const viaCepURL = "http://viacep.com.br/ws/%s/json/"

func NewClimateService(client HTTPClient) gateway.ZipCodeService {
	return &climateService{client: client}
}

type climateService struct {
	client HTTPClient
}

func (s *climateService) Get(ctx context.Context, zipCodeID model.ZipCodeID) (*model.ZipCode, error) {
	url := fmt.Sprintf(viaCepURL, zipCodeID)
	resp, err := s.client.Get(url)
	if err != nil {
		return &model.ZipCode{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &model.ZipCode{}, errors.New("failed to fetch data from ViaCEP")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &model.ZipCode{}, err
	}

	var zipCode model.ZipCode
	if err := json.Unmarshal(body, &zipCode); err != nil {
		return &model.ZipCode{}, err
	}

	if zipCode == (model.ZipCode{}) {
		return &model.ZipCode{}, errors.New(model.ErrZipCodeNotFound)
	}

	return &zipCode, nil
}
