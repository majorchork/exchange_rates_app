package service

import (
	"os"
	"testing"
	"time"

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
	exrate, err := rateService.GetLatestExchange()
	if err != nil {
		t.Error(err)
	}
	if len(exrate.Rates) == 0 {
		t.Error()
	}
	if exrate.Base == "" {
		t.Error()
	}
	dateString := exrate.Date
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		t.Error(err)
	}
	exrate, err = rateService.GetExchangeByDate(date)
	if err != nil {
		t.Error(err)
	}

	if len(exrate.Rates) == 0 {
		t.Error()
	}
	exrate, err = rateService.GetAnalyzedRates()
	if err != nil {
		t.Error(err)
	}

	if len(exrate.AnalyzedRates) == 0 {
		t.Error()
	}

	if err := os.Remove("test.db"); err != nil {
		t.Error(err)
	}
}
