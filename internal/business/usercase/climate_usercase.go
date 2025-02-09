package usercase

import (
	"climateZpCode/internal/business/gateway"
	"climateZpCode/internal/business/model"
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

type (
	ClimateUserCase interface {
		Get(ctx context.Context, zipCodeID model.ZipCodeID) (model.Climate, error)
	}

	climateUserCase struct {
		zipCodeService gateway.ZipCodeService
		weatherService gateway.WheaterService
		trace          trace.Tracer
	}
)

func NewClimateUserCase(zipCodeService gateway.ZipCodeService, weatherService gateway.WheaterService, trace trace.Tracer) ClimateUserCase {
	return climateUserCase{
		zipCodeService: zipCodeService,
		weatherService: weatherService,
		trace:          trace,
	}
}

func (u climateUserCase) Get(ctx context.Context, zipCodeID model.ZipCodeID) (model.Climate, error) {
	ctxSpan, span := u.trace.Start(ctx, "climate-user-case")

	traceID := span.SpanContext().TraceID().String()
	fmt.Println("User Case TraceID: ", traceID)

	if err := zipCodeID.Validate(); err != nil {
		return model.Climate{}, err
	}

	z, err := u.zipCodeService.Get(ctxSpan, zipCodeID)
	if err != nil {
		return model.Climate{}, err
	}

	weatherID := model.WeatherID(z.Localidade)
	weather, err := u.weatherService.Get(ctxSpan, weatherID)
	if err != nil {
		return model.Climate{}, err
	}

	var climate model.Climate

	climate.TempK = weather.Current.TempC + 273
	climate.TempC = weather.Current.TempC
	climate.TempF = weather.Current.TempC*1.8 + 32

	return climate, nil
}
