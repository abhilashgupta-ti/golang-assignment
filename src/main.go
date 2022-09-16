package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Rates struct {
	EUR string `json:"EUR"`
	USD string `json:"USD"`
}

type BitcoinRates struct {
	Currbitcoin Rates `json:"bitcoin"`
}

type EnclosingStruct struct {
	Data BitcoinRates `json:"data"`
}

var rates = Rates{EUR: "43,947.8947", USD: "49,822.2917"}
var BRates = BitcoinRates{Currbitcoin: rates}
var ESRates = EnclosingStruct{Data: BRates}

func main() {
	router := gin.Default()
	router.GET("/prices", getPrices)

	router.Run("localhost:9999")
}

// getPrices responds with the EUR and USD prices as JSON.
func getPrices(c *gin.Context) {
	fmt.Println(ESRates)
	c.IndentedJSON(http.StatusOK, ESRates)
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
