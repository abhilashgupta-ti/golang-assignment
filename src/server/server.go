package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

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

// lastAPI used to track the prices
var lastApiUsed string

// RSRates stores the values seen
var RSRates ResponseStruct

type PriceTracker interface {
	GetPriceFromTracker() (ResponseStruct, error)
	GetName() string
	GetURL() string
}

// updateExchangeRates updates RSRates the global variable
func updateExchangeRates(prices ResponseStruct, timeChannel chan time.Time) {
	RSRates = prices          //update the structure with the got value
	timeChannel <- time.Now() // return time of update
}

// GetPrices responds to the API call with the EUR and USD prices as JSON.
func GetPrices(c *gin.Context) {
	log.Printf("/prices requested by %v\n", c.ClientIP())
	var currentTime = time.Now()
	var currentValidity = lastUpdatedTime.Add(time.Second * time.Duration(ExpiryTime))
	// fmt.Println(currentValidity, currentTime)

	var PriceTrackerAPI = os.Getenv("PRICE_TRACKER")
	// if we don't have the exchange rates or have stale values -> get latest values
	if lastUpdatedTime.IsZero() || lastApiUsed != PriceTrackerAPI || currentTime.After(currentValidity) {
		// fmt.Println(lastUpdatedTime.IsZero())
		var priceTrackerInterface PriceTracker
		var curRates ResponseStruct
		var err error
		timeChannel := make(chan time.Time) // channel for concurrency control

		switch strings.ToLower(PriceTrackerAPI) {
		case "":
			log.Println("Tracker API not set! Defaulting to coindesk for now!")
			log.Println("You can set which API to extract rates from by setting the PRICE_TRACKER environment variable.")
			fallthrough
		case "coindesk": // In case PriceTrackerAPI isn't set in the environment
			log.Println("Rates extracted from Coindesk API")
			var cdApi = CoinDeskApi{}
			priceTrackerInterface = cdApi
		default:
			var returnError = "Internal Server Error. Please report issue at https://github.com/agtelus/golang-assignment"
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"returncode": http.StatusServiceUnavailable, "error": returnError})
			log.Panicf("PANIC!!! Environment variable 'PriceTrackerAPI' set to invalid value '%v'. Set it to one of {coindesk}!!!\n", PriceTrackerAPI)

		}

		curRates, err = priceTrackerInterface.GetPriceFromTracker() //get current rates
		if err != nil {
			log.Printf("Couldn't get rates from the %v API URL %v", priceTrackerInterface.GetName(), priceTrackerInterface.GetURL())
			log.Printf("Please try again or use another API.")
			var returnError = "Couldn't get the latest rates. Please report issue at https://github.com/agtelus/golang-assignment"
			c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"returncode": http.StatusServiceUnavailable, "error": returnError})
			return
		}
		go updateExchangeRates(curRates, timeChannel) //update local variable in background process
		lastUpdatedTime = <-timeChannel               //ensure synchronicity
		lastApiUsed = PriceTrackerAPI

	}

	c.IndentedJSON(http.StatusOK, RSRates) //return indented json response

}
