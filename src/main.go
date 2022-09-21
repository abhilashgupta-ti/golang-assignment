package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ExpiryTime in seconds
var ExpiryTime int

// main()
func main() {
	var err error
	// Figure out expiry Time and set it to global variable
	expiryTimeStr := os.Getenv("EXPIRY_TIME")
	if expiryTimeStr == "" { //ENV variable EXPIRY isn't set. Default to 0
		ExpiryTime = 0
	} else {
		ExpiryTime, err = strconv.Atoi(expiryTimeStr)
		if err != nil {
			log.Fatal("Invalid Expiry time. Please set the env variable \"EXPIRY\" with seconds in integer")
		}
	}
	log.Printf("Expiry time=%d seconds\n", ExpiryTime)

	//Initialise the gin router
	router := gin.Default()
	//getPrices handles API requests to `/prices`
	router.GET("/prices", GetPrices)
	// attach the router to http.Server and start the server on port 9999
	router.Run("localhost:9999")
}
