//go:generate mockery --dir=. --output=../../infrastructure/service/mock --name=WheaterService --structname=MockWheaterService --outpkg=mock --filename=weather_service.go --disable-version-string
package gateway

import (
	"climateZpCode/internal/business/model"
	"context"
)

type WheaterService interface {
	Get(context.Context, model.WeatherID) (*model.Weather, error)
}
