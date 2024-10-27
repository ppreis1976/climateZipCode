//go:generate mockery --dir=. --output=../../infrastructure/service/mock --name=ZipCodeService --structname=MockZipCodeService --outpkg=mock --filename=zipcode_service.go --disable-version-string
package gateway

import (
	"climateZpCode/internal/business/model"
	"context"
)

type ZipCodeService interface {
	Get(context.Context, model.ZipCodeID) (*model.ZipCode, error)
}
