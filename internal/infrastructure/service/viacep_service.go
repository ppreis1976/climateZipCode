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

	"go.opentelemetry.io/otel/trace"
)

const viaCepURL = "http://viacep.com.br/ws/%s/json/"

func NewClimateService(client HTTPClient, tracer trace.Tracer) gateway.ZipCodeService {
	return &climateService{
		client: client,
		tracer: tracer,
	}
}

type climateService struct {
	client HTTPClient
	tracer trace.Tracer
}

func (s *climateService) Get(ctx context.Context, zipCodeID model.ZipCodeID) (*model.ZipCode, error) {
	_, span := s.tracer.Start(ctx, "weather-api.get")
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	fmt.Println("Via CEP Api TraceID: ", traceID)

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
