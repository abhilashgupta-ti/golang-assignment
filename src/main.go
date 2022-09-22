package main

import (
	"github.com/agtelus/golang-assignment/src/server"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	// Get Bitcoin exchanges rates/prices
	router.GET("/prices", server.GetPrices)

	// Add more APIs here

	return router
}

// main()
func main() {
	// Figure out expiry Time and set it
	server.SetServerExpiryTime()

	//Initialise the gin router
	router := setupRouter()

	// attach the router to http.Server and start the server on port 9999
	router.Run("localhost:9999")
}
