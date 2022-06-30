package service

import (
	"encoding/xml"
	"log"
	"net/http"
	"time"

	"github.com/majorchork/rates_app/database"
	"github.com/majorchork/rates_app/models"
)

type ExchangeRateService interface {
	GetLatestExchange() (models.ExchangeRatesResponse, error)
	GetExchangeByDate(date time.Time) (models.ExchangeRatesResponse, error)
	GetAnalyzedRates() (models.ExchangeRatesResponse, error)
}

type Rates struct {
	BaseCurrency string
	DB           database.DbRepository
}

func NewRatesService() (ExchangeRateService, error) {
	// initialize the sql lite database
	SqliteDb := database.NewSqliteDb("ExchangeRate")

	// initialize the ratesService struct and pass in the repository
	ratesService := new(Rates)
	ratesService.BaseCurrency = "EUR"
	var err error
	ratesService.DB, err = database.NewRateRepo(SqliteDb)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ratesService.Prepare()

	return ratesService, nil
}

func (rh Rates) Prepare() {
	client := http.Client{}
	response, err := client.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml")
	if err != nil {
		log.Fatal(err)
	}

	var envelope models.Envelope
	if err = xml.NewDecoder(response.Body).Decode(&envelope); err != nil {
		log.Fatal(err)
	}
	log.Println("resetting database")
	if err := rh.DB.ResetDB(); err != nil {
		log.Fatal(err)
	}
	if err := rh.DB.Store(envelope); err != nil {
		log.Fatal(err)
	}
	log.Println("finished seeding")
}

// GetLatestExchange returns rates with DTO ForexResponse
func (rh Rates) GetLatestExchange() (ExchangeRateResponse models.ExchangeRatesResponse, err error) {
	result, err := rh.DB.FindByLatestDate()
	if err != nil {
		return models.ExchangeRatesResponse{}, err
	}

	ExchangeRateResponse.Base = rh.BaseCurrency
	ExchangeRateResponse.Date = result[0].ForexDate.Format("2006-01-02")
	ExchangeRateResponse.Rates = make(map[string]float64, 0)
	for _, resp := range result {
		ExchangeRateResponse.Rates[resp.Currency] = resp.Rate
	}

	return ExchangeRateResponse, nil
}

// GetExchangeByDate returns rates with certain date
func (rh Rates) GetExchangeByDate(date time.Time) (ExchangeRateResponse models.ExchangeRatesResponse, err error) {
	result, err := rh.DB.FindByDateString(date)
	if err != nil {
		return models.ExchangeRatesResponse{}, err
	}

	ExchangeRateResponse.Base = rh.BaseCurrency
	ExchangeRateResponse.Rates = make(map[string]float64, 0)
	for _, fx := range result {
		ExchangeRateResponse.Rates[fx.Currency] = fx.Rate
		ExchangeRateResponse.Date = fx.ForexDate.Format("2006-01-02")
	}

	return ExchangeRateResponse, nil
}

// GetAnalyzedRates returns analyzed rate: max, min, average
func (rh Rates) GetAnalyzedRates() (ExchangeRateResponse models.ExchangeRatesResponse, err error) {
	result, err := rh.DB.Find()
	if err != nil {
		return models.ExchangeRatesResponse{}, err
	}

	ExchangeRateResponse.Base = rh.BaseCurrency
	ExchangeRateResponse.AnalyzedRates = make(map[string]models.AnalyzedRates, 0)
	for _, fx := range result {
		rates, ok := ExchangeRateResponse.AnalyzedRates[fx.Currency]
		if !ok {
			ExchangeRateResponse.AnalyzedRates[fx.Currency] = models.AnalyzedRates{
				Min:   fx.Rate,
				Max:   fx.Rate,
				Sum:   fx.Rate,
				Count: 1,
			}
		} else {
			rates := ExchangeRateResponse.AnalyzedRates[fx.Currency]
			rates.Sum = rates.Sum + fx.Rate
			rates.Count = rates.Count + 1
			ExchangeRateResponse.AnalyzedRates[fx.Currency] = rates
		}

		if fx.Rate < rates.Min {
			rates := ExchangeRateResponse.AnalyzedRates[fx.Currency]
			rates.Min = fx.Rate
			ExchangeRateResponse.AnalyzedRates[fx.Currency] = rates
		}

		if fx.Rate > rates.Max {
			rates := ExchangeRateResponse.AnalyzedRates[fx.Currency]
			rates.Max = fx.Rate
			ExchangeRateResponse.AnalyzedRates[fx.Currency] = rates
		}
	}

	for currency := range ExchangeRateResponse.AnalyzedRates {
		rates := ExchangeRateResponse.AnalyzedRates[currency]
		rates.Avg = rates.Sum / float64(rates.Count)
		ExchangeRateResponse.AnalyzedRates[currency] = rates
	}

	return ExchangeRateResponse, nil
}
