package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/daragao/goUntitled/controllers"
	"github.com/gorilla/handlers"
)

func TestApp(t *testing.T) {
	//server
	ts := httptest.NewServer(handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateAppRouter()))
	defer ts.Close()

	_, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	//recorder
	/*
		t.Log("TestApp")
		req, err := http.NewRequest("GET", "http://localhost:3333/register", nil)
		if err != nil {
			log.Fatal(err)
		}

		w := httptest.NewRecorder()
		controllers.Register(w, req)
	*/
}
