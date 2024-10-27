package main

import (
	"climateZpCode/internal/infrastructure/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/climate", handler.ClimateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
