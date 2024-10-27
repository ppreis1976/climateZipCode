package handler

import (
	"climateZpCode/internal/business/model"
	"climateZpCode/internal/business/usercase"
	"climateZpCode/internal/infrastructure/service"
	"encoding/json"
	"errors"
	"net/http"
)

const ZipCode = "zip_code"

func ClimateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleZipCodeError(w, errors.New(model.ErrMethodNotAllowed))
		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	zipCode := r.Header.Get(ZipCode)
	if zipCode == "" {
		http.Error(w, "zip_code header is missing", http.StatusBadRequest)
		return
	}

	client := &service.DefaultHTTPClient{}

	zipCodeService := service.NewClimateService(client)
	weatherService := service.NewWeatherService(client)
	zipCodeUserCase := usercase.NewClimateUserCase(zipCodeService, weatherService)

	climate, err := zipCodeUserCase.Get(model.ZipCodeID(zipCode))
	if err != nil {
		handleZipCodeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(climate)
	if err != nil {
		return
	}
}

func handleZipCodeError(w http.ResponseWriter, err error) {
	if err.Error() == model.ErrMethodNotAllowed {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	if err.Error() == model.ErrZipCodeIDInvalid {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err.Error() == model.ErrZipCodeNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
