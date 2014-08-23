package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/daragao/goUntitled/config"
	"github.com/daragao/goUntitled/controllers"
	"github.com/daragao/goUntitled/models"
	"github.com/gorilla/handlers"
)

var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	// Setup the global variables and settings
	err := models.Setup()
	if err != nil {
		fmt.Println(err)
	}
	// Start the web server
	fmt.Printf("Starting app.go at %s\n", config.Conf.AppURL)
	http.ListenAndServe(config.Conf.AppURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateAppRouter()))
}
