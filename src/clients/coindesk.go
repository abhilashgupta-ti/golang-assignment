package clients

import (
	"encoding/json"
	"github.com/agtelus/golang-assignment/src/interfaces"
	"io/ioutil"
	"log"
	"net/http"
)

// Structs to extract the relevant information from the CoinDesk API
type RateStruct struct {
	Rate string `json:"rate"`
}

type BpiStruct struct {
	Usd RateStruct `json:"USD"`
	Eur RateStruct `json:"EUR"`
}

type RequestParsingStruct struct {
	Bpi BpiStruct `json:"bpi"`
}

type CoinDeskApi struct {
}

func (cdApi CoinDeskApi) GetName() string {
	return "CoinDesk"
}

func (cdApi CoinDeskApi) GetURL() string {
	return "https://api.coindesk.com/v1/bpi/currentprice.json"
}

// Get prices from CoinDesk API
func (cdApi CoinDeskApi) GetPriceFromTracker() (interfaces.ResponseStruct, error) {
	var updatedRates = interfaces.ResponseStruct{}
	var parsedRates = RequestParsingStruct{}
	requestURL := cdApi.GetURL()
	response, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Could not get latest rates from coindesk: %s\n", err)
		return updatedRates, err
	}
	responseJsonData, err := ioutil.ReadAll(response.Body) //convert response body into string
	if err != nil {
		log.Printf("couldn't read the Response body received from coindesk %s\n", err)
		return updatedRates, err
	}
	// Unmarshall the relevant part of the json into our structure
	json.Unmarshal([]byte(responseJsonData), &parsedRates)
	var BtcRates = interfaces.BitcoinRatesStruct{parsedRates.Bpi.Eur.Rate, parsedRates.Bpi.Usd.Rate}
	updatedRates = interfaces.ResponseStruct{Data: interfaces.DataStruct{Bitcoin: BtcRates}}
	return updatedRates, nil
}
