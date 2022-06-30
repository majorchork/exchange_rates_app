package server

import (
	"encoding/json"
	"errors"
	"github.com/majorchork/rates_app/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockService "github.com/majorchork/rates_app/service/mocks"
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
	mockRateService := mockService.NewMockExchangeRateService(ctrl)
	h := &handler{
		rateService: mockRateService,
	}

	body, err := json.Marshal(expected)
	if err != nil {
		t.Error(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/rates/latest", h.GetLatestExchangeHandler())
	t.Run("Testing for error", func(t *testing.T) {

		mockRateService.EXPECT().GetLatestExchange().Return(models.ExchangeRatesResponse{}, errors.New("error"))
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/rates/latest", nil)
		if err != nil {
			t.Error(err)
		}
		mux.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "error")

	})
	t.Run("Testing for Successful Request", func(t *testing.T) {

		mockRateService.EXPECT().GetLatestExchange().Return(expected, nil)
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/rates/latest", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("Content-Type", "application/json")
		mux.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), string(body))

	})

}

func TestGetExchangeByDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRateService := mockService.NewMockExchangeRateService(ctrl)
	h := NewHandler(mockRateService)

	expected := models.ExchangeRatesResponse{
		Base: "TEST",
		Date: "2006-06-24",
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
	mux.HandleFunc("/rates/", h.GetExchangeByDateHandler())
	t.Run("Testing for bad Request", func(t *testing.T) {

		mockRateService.EXPECT().GetExchangeByDate(time.Date(2006, 06, 24, 0, 0, 0, 0, time.UTC)).Return(models.ExchangeRatesResponse{}, errors.New("error"))
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/rates/2006-06-24", nil)
		if err != nil {
			t.Error(err)
		}
		mux.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "error")

	})
	t.Run("Testing for Successful Request", func(t *testing.T) {

		mockRateService.EXPECT().GetExchangeByDate(time.Date(2006, 06, 24, 0, 0, 0, 0, time.UTC)).Return(expected, nil)
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/rates/2006-06-24", nil)
		if err != nil {
			t.Error(err)
		}

		mux.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), string(body))

	})

}

func TestGetAnalyzedRates(t *testing.T) {
	expected := models.ExchangeRatesResponse{
		Base: "TEST",
		Date: "2022-06-24",
		AnalyzedRates: map[string]models.AnalyzedRates{
			"AUD": {
				1.4994,
				1.5693,
				1.5340524590163933,
				0,
				0,
			},
			"BGN": {
				1.9558,
				1.9558,
				1.9557999999999973,
				0,
				0,
			},
			"USD": {
				1.1562,
				1.2065,
				1.1783852459016388,
				0,
				0,
			},
			"ZAR": {
				14.7325,
				17.0212,
				16.06074426229508,
				0,
				0,
			},
		},
	}
	ctrl := gomock.NewController(t)
	mockRateService := mockService.NewMockExchangeRateService(ctrl)
	h := NewHandler(mockRateService)

	mux := http.NewServeMux()
	mux.HandleFunc("/rates/analyze", h.GetAnalyzedRatesHandler())
	body, err := json.Marshal(expected)
	if err != nil {
		t.Error(err)
	}
	t.Run("Testing for error", func(t *testing.T) {

		mockRateService.EXPECT().GetAnalyzedRates().Return(models.ExchangeRatesResponse{}, errors.New("error"))
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/rates/analyze", nil)
		if err != nil {
			t.Error(err)
		}

		mux.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "error")

	})
	t.Run("Testing for Successful Request", func(t *testing.T) {

		mockRateService.EXPECT().GetAnalyzedRates().Return(expected, nil)
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/rates/analyze", nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("Content-Type", "application/json")
		mux.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), string(body))

	})

}
