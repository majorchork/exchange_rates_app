package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
			respond(w, nil, 500)
		}
		respond(w, resp, 200)
	}
}

func (h *handler) GetExchangeByDateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		date, ok := query["date"]

		if !ok || len(date) == 0 {
			// TODO: return a proper error message to the user
			log.Fatal("something went wrong")
		}

		dateString := strings.Join(date, "-")
		resp, err := h.rateService.GetExchangeByDate(dateString)
		if err != nil {
			respond(w, nil, 500)
		}

		respond(w, resp, 200)
	}
}

func (h *handler) GetAnalyzedRatesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := h.rateService.GetAnalyzedRates()
		if err != nil {
			respond(w, nil, 500)
		}

		respond(w, resp, 200)
	}

}
