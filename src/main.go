package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

// Structs to construct the JSON API response
type BitcoinRatesStruct struct {
	Eur string `json:"EUR"`
	Usd string `json:"USD"`
}

type DataStruct struct {
	Bitcoin BitcoinRatesStruct `json:"bitcoin"`
}

type ResponseStruct struct {
	Data DataStruct `json:"data"`
}

// global variables

// Stores the time of the last update of the values
var lastUpdatedTime time.Time

// Stores the values seen
var RSRates ResponseStruct

// The expiry time in seconds
var expiryTime int

// main()
func main() {
	var err error
	// Figure out expiry Time and set it to global variable
	expiryTimeStr := os.Getenv("EXPIRY")
	if expiryTimeStr == "" { //ENV variable EXPIRY isn't set. Default to 0
		expiryTime = 0
	} else {
		expiryTime, err = strconv.Atoi(expiryTimeStr)
		if err != nil {
			log.Fatal("Invalid Expiry time. Please set the env variable \"EXPIRY\" with seconds in integer")
		}
	}
	log.Printf("Expiry time=%d\n", expiryTime)

	//Initialise the gin router
	router := gin.Default()
	//getPrices handles API requests to `/prices`
	router.GET("/prices", getPrices)
	// attach the router to http.Server and start the server on port 9999
	router.Run("localhost:9999")
}

// updateExchangeRates write the  RSRatesthe global variable
func updateExchangeRates(BtcRates BitcoinRatesStruct, timeChannel chan time.Time) {
	RSRates = ResponseStruct{Data: DataStruct{Bitcoin: BtcRates}} //update the structure with the got value
	timeChannel <- time.Now()                                     // return time of update
}

// getPrices responds to the API call with the EUR and USD prices as JSON.
func getPrices(c *gin.Context) {
	log.Printf("/prices requested by %v\n", c.ClientIP())
	var currentTime = time.Now()
	var currentValidity = lastUpdatedTime.Add(time.Second * time.Duration(expiryTime))
	// fmt.Println(currentValidity, currentTime)

	// if we don't have the exchange rates or have stale values -> get latest values
	if lastUpdatedTime.IsZero() || currentTime.After(currentValidity) {
		// fmt.Println(lastUpdatedTime.IsZero())
		var curRates RequestParsingStruct
		var err error
		timeChannel := make(chan time.Time)        // channel for concurrency control
		curRates, err = getPricesFromCoinDeskAPI() //get current rates
		if err != nil {
			log.Fatal("Couldn't get rates")
		}
		var BtcRates = BitcoinRatesStruct{Eur: curRates.Bpi.Eur.Rate, Usd: curRates.Bpi.Usd.Rate}
		go updateExchangeRates(BtcRates, timeChannel) //update local variable in background process
		lastUpdatedTime = <-timeChannel               //ensure synchronicity
	}

	c.IndentedJSON(http.StatusOK, RSRates) //return indented json response
}

// Get prices from CoinDesk API
func getPricesFromCoinDeskAPI() (RequestParsingStruct, error) {
	var currentRates = RequestParsingStruct{}
	requestURL := "https://api.coindesk.com/v1/bpi/currentprice.json"
	response, err := http.Get(requestURL)
	if err != nil {
		log.Printf("client: could not get current rates from coindesk: %s\n", err)
		return currentRates, err
	}
	responseData, err := ioutil.ReadAll(response.Body) //convert response body into string
	if err != nil {
		log.Printf("couldn't read the Response body")
		return currentRates, err
	}
	returnedJson := responseData
	// Unmarshall the relevant part of the json into our structure
	json.Unmarshal([]byte(returnedJson), &currentRates)
	return currentRates, nil

}
