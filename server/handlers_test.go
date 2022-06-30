package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/majorchork/rates_app/models"
	mock_service "github.com/majorchork/rates_app/service/mocks"
)

func TestGetLatestExchange(t *testing.T) {
	expected := models.ExchangeRatesResponse{
		Base: "TEST",
		Date: "2006-01-02",
		Rates: map[string]float64{
			"AUD": 1.5339,
			"BGN": 1.9558,
			"USD": 1.2023,
			"ZAR": 14.8845,
		},
	}
	ctrl := gomock.NewController(t)
	mockRateService := mock_service.NewMockExchangeRateService(ctrl)
	_ = &handler{
		rateService: mockRateService,
	}
	//a := http.Server{Handler: h}
	//handlers := NewHandler(*mockHelper)
	body, err := json.Marshal(expected)
	if err != nil {
		t.Error(err)
	}

	mux := http.NewServeMux()
	t.Run("Testing for error", func(t *testing.T) {

		mockRateService.EXPECT().GetLatestExchange().Return(models.ExchangeRatesResponse{}, errors.New("error"))
		rw := httptest.NewRecorder()
		fmt.Println(rw.Code)
		req, err := http.NewRequest(http.MethodGet, "/rates/getlatest", strings.NewReader(string(body)))
		if err != nil {
			t.Error(err)
		}
		fmt.Println(rw.Code)
		mux.ServeHTTP(rw, req)
		//req.Header.Add("Content-Type", "application/json")
		//fmt.Println(rw.Body.String())
		fmt.Println(rw.Code)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "error")

	})
	t.Run("Testing for Successful Request", func(t *testing.T) {

		mockRateService.EXPECT().GetLatestExchange().Return(expected, nil)
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/rates/getlatest", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("Content-Type", "application/json")
		mux.ServeHTTP(rw, req)
		fmt.Println(rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), string(body))

	})

}

func TestGetExchangeByDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRateService := mock_service.NewMockExchangeRateService(ctrl)
	_ = NewHandler(mockRateService)

	expected := models.ExchangeRatesResponse{
		Base: "TEST",
		Date: "2006-01-02",
		Rates: map[string]float64{
			"AUD": 1.5339,
			"BGN": 1.9558,
			"USD": 1.2023,
			"ZAR": 14.8845,
		},
	}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Error(err)
	}

	mux := http.NewServeMux()
	t.Run("Testing for b Request", func(t *testing.T) {

		mockRateService.EXPECT().GetExchangeByDate("2006-01-02").Return(models.ExchangeRatesResponse{}, errors.New("error"))
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/rates?date=2006-01-02", nil)
		if err != nil {
			t.Error(err)
		}
		mux.ServeHTTP(rw, req)
		//req.Header.Add("Content-Type", "application/json")
		//fmt.Println(rw.Body.String())
		fmt.Println(rw.Code)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "error")

	})
	t.Run("Testing for Successful Request", func(t *testing.T) {

		mockRateService.EXPECT().GetExchangeByDate("2006-01-02").Return(expected, nil)
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/rates?date=2006-01-02", nil)
		if err != nil {
			t.Error(err)
		}
		mux.ServeHTTP(rw, req)
		//req.Header.Add("Content-Type", "application/json")
		fmt.Println(rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), string(body))

	})

}

//func TestGetAnalyzedRates(t *testing.T) {
//	expected := models.ExchangeRatesResponse{
//		Base: "TEST",
//		Date: "2022-06-24",
//		AnalyzedRates: map[string]models.AnalyzedRates{
//			"AUD": {
//				1.4994,
//				1.5693,
//				1.5340524590163933,
//				0,
//				0,
//			},
//			"BGN": {
//				1.9558,
//				1.9558,
//				1.9557999999999973,
//				0,
//				0,
//			},
//			"USD": {
//				1.1562,
//				1.2065,
//				1.1783852459016388,
//				0,
//				0,
//			},
//			"ZAR": {
//				14.7325,
//				17.0212,
//				16.06074426229508,
//				0,
//				0,
//			},
//		},
//	}
//	ctrl := gomock.NewController(t)
//	mockRateService := mock_service.NewMockExchangeRateService(ctrl)
//	_ = NewHandler(mockRateService)
//
//	t.Run("Testing for Successful Request", func(t *testing.T) {
//
//		mockRateService.EXPECT().GetAnalyzedRates().Return(expected)
//		rw := httptest.NewRecorder()
//		req, err := http.NewRequest(http.MethodGet, "/rates/analyze", nil)
//		if err != nil {
//			t.Error(err)
//		}
//		req.Header.Add("Content-Type", "application/json")
//		fmt.Println(rw.Body.String())
//
//	})
//
//}
