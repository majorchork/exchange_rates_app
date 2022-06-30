package server

import (
	"net/http"
)

func SetupRouter(h *handler) error {
	//router := http.NewServeMux()

	http.HandleFunc("/rates/latest", h.GetLatestExchangeHandler())
	http.HandleFunc("/rates/analyze", h.GetAnalyzedRatesHandler())
	http.HandleFunc("/rates/", h.GetExchangeByDateHandler())
	if err := http.ListenAndServe(":8085", nil); err != nil {
		return err
	}
	return nil
}
