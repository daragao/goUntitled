package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	//"github.com/daragao/goUntitled/auth"
	"github.com/daragao/goUntitled/controllers"
	"github.com/daragao/goUntitled/models"
	"github.com/gorilla/handlers"
)

func TestApp(t *testing.T) {
	//server

	err := models.Setup()
	if err != nil {
		fmt.Println(err)
	}

	ts := httptest.NewServer(handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateAppRouter()))
	defer ts.Close()

	jsonStr := "{ \"username\": \"test\", \"password\": \"ipass\"}"
	_, err = http.Post(ts.URL+"/register", "json/text", bytes.NewBufferString(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	//w := httptest.NewRecorder()
	//controllers.Register(w, req)

	//recorder
	/*t.Log("TestApp")

	jsonStr := "{ \"username\": \"test\", \"password\": \"ipass\"}"
	req, err := http.NewRequest("GET", "http://localhost:3333/register", bytes.NewBufferString(jsonStr))
	if err != nil {
		log.Fatal(err)
	}

	auth.Register(req)

	w := httptest.NewRecorder()
	controllers.Register(w, req)*/
}
