package models

import "time"

type Exchange struct {
	ID           int       `json:"id"`
	BaseCurrency string    `json:"-"`
	Currency     string    `json:"currency" validate:"required"`
	ForexDate    time.Time `json:"forex_date" validate:"required"`
	Rate         float64   `json:"rate" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
}
