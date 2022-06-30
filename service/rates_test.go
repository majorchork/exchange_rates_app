package service

import (
	"os"
	"testing"

	"github.com/majorchork/rates_app/database"
)

func TestRateService(t *testing.T) {
	sqlite := database.NewSqliteDb("test.db")

	rateService := new(Rates)
	rateService.BaseCurrency = "EUR"
	var err error
	rateService.DB, err = database.NewRateRepo(sqlite)
	if err != nil {
		t.Error(err)
	}
	rateService.Prepare()
	exrate := rateService.GetLatestExchange()
	if len(exrate.Rates) == 0 {
		t.Error()
	}
	if exrate.Base == "" {
		t.Error()
	}
	date := exrate.Date
	exrate = rateService.GetExchangeByDate(date)
	if len(exrate.Rates) == 0 {
		t.Error()
	}
	exrate = rateService.GetAnalyzedRates()
	if len(exrate.AnalyzedRates) == 0 {
		t.Error()
	}

	if err := os.Remove("test.db"); err != nil {
		t.Error(err)
	}
}
