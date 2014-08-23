package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/daragao/goUntitled/auth"
	mid "github.com/daragao/goUntitled/middleware"
	"github.com/daragao/goUntitled/models"
	ctx "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/justinas/nosurf"
)

var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func CreateAppRouter() http.Handler {
	Logger.Println("Route started")
	router := mux.NewRouter()
	// Base Front-end routes
	router.HandleFunc("/login", Login)
	router.HandleFunc("/logout", Use(Logout, mid.RequireLogin))
	router.HandleFunc("/register", Register)
	router.HandleFunc("/", Use(Base, mid.RequireLogin))

	csrfHandler := nosurf.New(router)
	csrfHandler.ExemptGlob("/register")

	return Use(csrfHandler.ServeHTTP, mid.GetContext)
}

// Use allows us to stack middleware to process the request
// Example taken from https://github.com/gorilla/mux/pull/36#issuecomment-25849172
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}

func Login(w http.ResponseWriter, r *http.Request) {
	Logger.Println("Login!!!")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	Logger.Println("Logout!!!")
}

func Register(w http.ResponseWriter, r *http.Request) {
	Logger.Println("Register!!!")
	// If it is a post request, attempt to register the account
	// Now that we are all registered, we can log the user in
	params := struct {
		Title   string
		Flashes []interface{}
		User    models.User
		Token   string
	}{Title: "Register", Token: nosurf.Token(r)}
	session := ctx.Get(r, "session").(*sessions.Session)
	switch {
	case r.Method == "GET":
		params.Flashes = session.Flashes()
		session.Save(r, w)
		//getTemplate(w, "register").ExecuteTemplate(w, "base", params)
	case r.Method == "POST":
		//Attempt to register
		succ, err := auth.Register(r)
		//If we've registered, redirect to the login page
		if succ {
			session.AddFlash(models.Flash{
				Type:    "success",
				Message: "Registration successful!.",
			})
			session.Save(r, w)
			//http.Redirect(w, r, "/login", 302)
		} else {
			// Check the error
			m := ""
			if err == models.ErrUsernameTaken {
				m = "Username already taken"
			} else {
				m = "Unknown error - please try again"
				Logger.Println(err)
			}
			session.AddFlash(models.Flash{
				Type:    "danger",
				Message: m,
			})
			session.Save(r, w)
			//http.Redirect(w, r, "/register", 302)
		}
	}
}

func Base(w http.ResponseWriter, r *http.Request) {
	Logger.Println("Base!!!")
}
