package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/majorchork/rates_app/service"
)

type handler struct {
	rateService service.ExchangeRateService
}

func NewHandler(ratesService service.ExchangeRateService) *handler {
	return &handler{rateService: ratesService}
}

func respond(w http.ResponseWriter, data interface{}, statuscode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("error marshaling the response: %v", err)
		return
	}
	w.Write(response)
}

func (h *handler) GetLatestExchangeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := h.rateService.GetLatestExchange()
		if err != nil {
			respond(w, "error", http.StatusInternalServerError)
		}
		respond(w, resp, http.StatusOK)
	}
}

func (h *handler) GetExchangeByDateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateString := strings.TrimPrefix(r.URL.Path, "/rates/")
		date, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			respond(w, "invalid data", http.StatusBadRequest)
		}

		//dateString := strings.Join(date, "-")
		resp, err := h.rateService.GetExchangeByDate(date)
		if err != nil {
			respond(w, "error", http.StatusInternalServerError)
		}

		respond(w, resp, http.StatusOK)
	}
}

func (h *handler) GetAnalyzedRatesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := h.rateService.GetAnalyzedRates()
		if err != nil {
			respond(w, "error", http.StatusInternalServerError)
		}

		respond(w, resp, http.StatusOK)
	}

}
