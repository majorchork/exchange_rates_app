package server

import (
	"github.com/majorchork/rates_app/service"
)

func Run() error {
	ratesService, err := service.NewRatesService()
	if err != nil {
		return err
	}

	return SetupRouter(NewHandler(ratesService))
}
