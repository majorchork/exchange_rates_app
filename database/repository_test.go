package database

import (
	"encoding/xml"
	"github.com/majorchork/rates_app/models"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestDb_Interface(t *testing.T) {
	sqlite := NewSqliteDb("test.db")

	dbHandler, err := NewRateRepo(sqlite)
	if err != nil {
		t.Error(err)
	}

	if err = dbHandler.ResetDB(); err != nil {
		t.Error(err)
	}

	xmlFile, err := os.Open("test.xml")
	if err != nil {
		t.Error(err)
	}
	var envelope models.Envelope
	data, err := ioutil.ReadAll(xmlFile)
	if err := xml.Unmarshal(data, &envelope); err != nil {
		t.Error(err)
	}
	if err := xmlFile.Close(); err != nil {
		t.Error(err)
	}
	if err = dbHandler.Store(envelope); err != nil {
		t.Error(err)
	}
	rate, err := dbHandler.FindByLatestDate()
	if err != nil {
		t.Error(err)
	}
	if len(rate) == 0 {
		t.Error()
	}
	rate, err = dbHandler.FindByDateString(time.Date(2022, 06, 24, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Error(err)
	}
	if len(rate) == 0 {
		t.Error()
	}

	rate, err = dbHandler.Find()
	if err != nil {
		t.Error(err)
	}
	if len(rate) == 0 {
		t.Error()
	}
	if err := os.Remove("test.db"); err != nil {
		t.Error(err)
	}
}
