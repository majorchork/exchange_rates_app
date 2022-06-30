package database

import (
	"fmt"
	"github.com/majorchork/rates_app/models"
	"time"
)

type DbRepository interface {
	ResetDB() error
	Store(envelope models.Envelope) error
	FindByLatestDate() ([]models.Exchange, error)
	FindByDateString(date time.Time) ([]models.Exchange, error)
	Find() ([]models.Exchange, error)
}

// DBHandler a database interface handler that has 2 methods
type DBHandler interface {
	Execute(statement string) error
	Query(statement string) Row
}

// Row executes 2 methods scan & next
type Row interface {
	Scan(dest ...interface{}) error
	Next() bool
}

type RateRepo struct {
	DB DBHandler
}

func NewRateRepo(sqlite DBHandler) (*RateRepo, error) {
	dbExchange := new(RateRepo)
	//dbExchange.dbHandlers = sqlite
	dbExchange.DB = sqlite
	return dbExchange, nil
}

func (repo *RateRepo) ResetDB() error {
	if err := repo.DB.Execute("DROP TABLE IF EXISTS ExchangeRate;"); err != nil {
		return err
	}

	if err := repo.DB.Execute(`
		CREATE TABLE ExchangeRate
		(   
			id INTEGER CONSTRAINT exchange_pk PRIMARY KEY AUTOINCREMENT,
			currency TEXT NOT NULL,
			forex_date TIMESTAMP NOT NULL,
			rate FLOAT,
			createdAt DATE DEFAULT CURRENT_TIMESTAMP NOT NULL
		);
	`); err != nil {
		return err
	}

	return nil
}

// Store data from ECB rates to database
func (repo *RateRepo) Store(envelope models.Envelope) error {
	for _, currencies := range envelope.Exchanges.CurrenciesPerDate {
		for _, currency := range currencies.Currency {
			if err := repo.DB.Execute(fmt.Sprintf("INSERT INTO ExchangeRate (currency, rate, forex_date) VALUES ('%s', '%s', '%s')", currency.Currency, currency.Rate, currencies.Time)); err != nil {
				return err
			}
		}
	}

	return nil
}

// FindByLatestDate returns exchanges with latest date from database
func (repo *RateRepo) FindByLatestDate() ([]models.Exchange, error) {
	row := repo.DB.Query(`
		SELECT currency, rate, forex_date from ExchangeRate  
WHERE forex_date = (SELECT  max(forex_date) FROM ExchangeRate ) 
ORDER by forex_date asc;
	`)
	var exchanges []models.Exchange
	for row.Next() {
		var exchange models.Exchange
		if err := row.Scan(&exchange.Currency, &exchange.Rate, &exchange.ForexDate); err != nil {
			return nil, err
		}
		exchanges = append(exchanges, exchange)
	}
	return exchanges, nil
}

// FindByDateString returns exchanges
func (repo *RateRepo) FindByDateString(date time.Time) ([]models.Exchange, error) {
	row := repo.DB.Query(fmt.Sprintf(`
		SELECT t.currency, t.rate, t.forex_date
		FROM ExchangeRate t
		WHERE t.forex_date = '%s'
		ORDER BY t.currency;
	`, date.Format("2006-01-02")))
	var exchanges []models.Exchange
	for row.Next() {
		var exchange models.Exchange
		if err := row.Scan(&exchange.Currency, &exchange.Rate, &exchange.ForexDate); err != nil {
			return nil, err
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges, nil
}

// Find returns all exchanges
func (repo *RateRepo) Find() ([]models.Exchange, error) {
	row := repo.DB.Query(`
		SELECT currency, rate, forex_date
		FROM ExchangeRate
		ORDER BY currency;
	`)
	var exchanges []models.Exchange
	for row.Next() {
		var exchange models.Exchange
		if err := row.Scan(&exchange.Currency, &exchange.Rate, &exchange.ForexDate); err != nil {
			return nil, err
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges, nil
}
