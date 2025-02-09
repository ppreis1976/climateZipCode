package handler

import (
	"climateZpCode/internal/business/model"
	"climateZpCode/internal/business/usercase"
	"climateZpCode/internal/infrastructure/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	RequestID = "X-Request-ID"
	ZipCode   = "zip_code"
)

func ClimateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	if r.Method != http.MethodGet {
		handleZipCodeError(w, errors.New(model.ErrMethodNotAllowed))
		return
	}

	tracer := service.NewTracer()

	ctxSpan, span := tracer.Start(ctx, "climate-api-handler")
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	fmt.Println("handler TraceID: ", traceID)

	trace.SpanFromContext(ctxSpan).AddEvent("climate handler before usercase")

	zipCode := r.Header.Get(ZipCode)
	if zipCode == "" {
		http.Error(w, "zip_code header is missing", http.StatusBadRequest)
		return
	}

	client := &service.DefaultHTTPClient{}

	zipCodeService := service.NewClimateService(client, tracer)
	weatherService := service.NewWeatherService(client, tracer)
	zipCodeUserCase := usercase.NewClimateUserCase(zipCodeService, weatherService, tracer)

	climate, err := zipCodeUserCase.Get(ctxSpan, model.ZipCodeID(zipCode))
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
