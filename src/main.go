package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

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

// getPrices responds with the EUR and USD prices as JSON.
func getPrices(c *gin.Context) {
	// fmt.Println(RSRates)
	curRates, err := getPricesFromCoinBaseAPI()
	if err != nil {
		log.Fatal("Couldn't get rates")
	}
	var BtcRates = BitcoinRatesStruct{Eur: curRates.Bpi.Eur.Rate, Usd: curRates.Bpi.Usd.Rate}
	RSRates := ResponseStruct{Data: DataStruct{Bitcoin: BtcRates}}
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

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// )

// func main() {
// 	data1 := map[string]interface{}{
// 		"data": map[string]interface{}{
// 			"bitcoin": map[string]interface{}{
// 				"EUR": "43,947.8947",
// 				"USD": "49,822.2917",
// 			},
// 		},
// 	}
// 	// data := map[string]interface{}{
// 	// 	"intValue":    1234,
// 	// 	"boolValue":   true,
// 	// 	"stringValue": "hello!",
// 	// 	"objectValue": map[string]interface{}{
// 	// 		"arrayValue": []int{1, 2, 3, 4},
// 	// 	},
// 	// }

// 	jsonData, err := json.MarshalIndent(data1, "", "   ")
// 	if err != nil {
// 		fmt.Printf("could not marshal json: %s\n", err)
// 		return
// 	}

// 	fmt.Printf("%s\n", jsonData)
// }
