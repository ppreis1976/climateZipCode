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
	"strings"

	"go.opentelemetry.io/otel/trace"
)

const weatherURL = "http://api.weatherapi.com/v1/current.json?key=6ce46037f7034601a4e231757240909&q=%s&aqi=no"

func NewWeatherService(client HTTPClient, tracer trace.Tracer) gateway.WheaterService {
	return &weatherService{
		client: client,
		tracer: tracer,
	}
}

type weatherService struct {
	client HTTPClient
	tracer trace.Tracer
}

func (s *weatherService) Get(ctx context.Context, weatherID model.WeatherID) (*model.Weather, error) {
	_, span := s.tracer.Start(ctx, "weather-api.get")
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	fmt.Println("Weather Api TraceID: ", traceID)

	url := fmt.Sprintf(weatherURL, encodeSpaces(removeAccents(string(weatherID))))
	resp, err := s.client.Get(url)
	if err != nil {
		return &model.Weather{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &model.Weather{}, errors.New("failed to fetch data from weatherapi")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &model.Weather{}, err
	}

	var weather model.Weather
	if err := json.Unmarshal(body, &weather); err != nil {
		return &model.Weather{}, err
	}

	if weather == (model.Weather{}) {
		return &model.Weather{}, errors.New(model.ErrZipCodeNotFound)
	}

	return &weather, nil
}

func removeAccents(input string) string {
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "ã", "a", "â", "a", "ä", "a",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"í", "i", "ì", "i", "î", "i", "ï", "i",
		"ó", "o", "ò", "o", "õ", "o", "ô", "o", "ö", "o",
		"ú", "u", "ù", "u", "û", "u", "ü", "u",
		"ç", "c",
		"Á", "A", "À", "A", "Ã", "A", "Â", "A", "Ä", "A",
		"É", "E", "È", "E", "Ê", "E", "Ë", "E",
		"Í", "I", "Ì", "I", "Î", "I", "Ï", "I",
		"Ó", "O", "Ò", "O", "Õ", "O", "Ô", "O", "Ö", "O",
		"Ú", "U", "Ù", "U", "Û", "U", "Ü", "U",
		"Ç", "C",
	)
	return replacer.Replace(input)
}

func encodeSpaces(input string) string {
	return strings.ReplaceAll(input, " ", "%20")
}
