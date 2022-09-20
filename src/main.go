package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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

func main() {
	router := gin.Default()
	router.GET("/prices", getPrices)

	router.Run("localhost:9999")
}

var lastUpdatedTime time.Time
var RSRates ResponseStruct

// getPrices responds with the EUR and USD prices as JSON.
func getPrices(c *gin.Context) {
	// fmt.Println(RSRates)
	var validTime int
	validTimeStr := os.Getenv("EXPIRY")
	var err error

	fmt.Println(validTimeStr)

	if validTimeStr == "" {
		validTime = 0
	} else {
		validTime, err = strconv.Atoi(validTimeStr)
		if err != nil {
			log.Fatal("Invalid Expiry time. Please set the env variable \"EXPIRY\" with seconds in integer")
		}
	}
	var currentTime = time.Now()
	fmt.Println(lastUpdatedTime.Add(time.Second*time.Duration(validTime)), currentTime)
	if lastUpdatedTime.IsZero() || currentTime.After(lastUpdatedTime.Add(time.Second*time.Duration(validTime))) {
		fmt.Println(lastUpdatedTime.IsZero())

		lastUpdatedTime = currentTime
		var curRates RequestParsingStruct
		curRates, err = getPricesFromCoinBaseAPI()
		if err != nil {
			log.Fatal("Couldn't get rates")
		}
		var BtcRates = BitcoinRatesStruct{Eur: curRates.Bpi.Eur.Rate, Usd: curRates.Bpi.Usd.Rate}
		RSRates = ResponseStruct{Data: DataStruct{Bitcoin: BtcRates}}
	}

	c.IndentedJSON(http.StatusOK, RSRates)
}

func getPricesFromCoinBaseAPI() (RequestParsingStruct, error) {
	var currentRates = RequestParsingStruct{}
	requestURL := "https://api.coindesk.com/v1/bpi/currentprice.json"
	response, err := http.Get(requestURL)
	if err != nil {
		log.Printf("client: could not get current rates: %s\n", err)
		return currentRates, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("couldn't read the Response body")
		return currentRates, err
	}
	// fmt.Println(responseData)
	returnedJson := responseData

	json.Unmarshal([]byte(returnedJson), &currentRates)
	return currentRates, nil

}
