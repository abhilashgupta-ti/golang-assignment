package server

import (
	"github.com/agtelus/golang-assignment/src/clients"
	"github.com/agtelus/golang-assignment/src/interfaces"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// global variables

// ExpiryTime in seconds
var ExpiryTime int

// Environment variable to set the API from which we gather our exchange rates
var PriceTrackerAPI string

// Stores the time of the last update of the values
var lastUpdatedTime time.Time

// lastAPI used to track the prices
var lastApiUsed string

// RSRates stores the latest price values
var RSRates interfaces.ResponseStruct

var coinDeskApi = clients.CoinDeskApi{}

// updateExchangeRates updates RSRates the global variable
func updateExchangeRates(prices interfaces.ResponseStruct, timeChannel chan time.Time) {
	RSRates = prices          //update the structure with the got value
	timeChannel <- time.Now() // return time of update
}

// GetPrices responds to the API call with the EUR and USD prices as JSON.
func GetPrices(c *gin.Context) {
	log.Printf("/prices requested by %v\n", c.ClientIP())
	var currentTime = time.Now()
	var currentValidity = lastUpdatedTime.Add(time.Second * time.Duration(ExpiryTime))

	var fetchFreshData = lastUpdatedTime.IsZero() || lastApiUsed != PriceTrackerAPI || currentTime.After(currentValidity)
	if fetchFreshData {
		var priceTrackerInterface interfaces.PriceTracker
		var curRates interfaces.ResponseStruct
		var err error
		timeChannel := make(chan time.Time) // channel for concurrency control

		switch PriceTrackerAPI {
		case "":
			log.Println("Tracker API not set! Defaulting to coindesk for now!")
			log.Println("You can set which API to extract rates from by setting the PRICE_TRACKER environment variable.")
			fallthrough
		case "coindesk": // In case PriceTrackerAPI isn't set in the environment
			log.Println("Rates extracted from Coindesk API")
			priceTrackerInterface = coinDeskApi
		default:
			var returnError = "Internal Server Error. Please report issue at https://github.com/agtelus/golang-assignment"
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"returncode": http.StatusServiceUnavailable, "error": returnError})
			log.Panicf("PANIC!!! Environment variable 'PRICE_TRACKER' set to invalid value '%v'. Set it to one of {coindesk}!!!\n", PriceTrackerAPI)

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

func SetServerExpiryTime() {
	var err error
	expiryTimeStr := os.Getenv("EXPIRY_TIME")
	if expiryTimeStr == "" { //ENV variable EXPIRY isn't set. Default to 0
		ExpiryTime = 0
	} else {
		ExpiryTime, err = strconv.Atoi(expiryTimeStr)
		if err != nil {
			log.Fatal("Invalid Expiry time. Please set the env variable \"EXPIRY\" with seconds in integer")
		}
	}
	log.Printf("Expiry time set to %d seconds\n", ExpiryTime)
}

func SetPriceTrackerAPI() {
	PriceTrackerAPI = strings.ToLower(os.Getenv("PRICE_TRACKER"))
}
