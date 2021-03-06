package auth

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"log"
	"os"

	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"github.com/daragao/goUntitled/models"
	ctx "github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

//init registers the necessary models to be saved in the session later
func init() {
	gob.Register(&models.User{})
	gob.Register(&models.Flash{})
}

var Store = sessions.NewCookieStore(
	[]byte(securecookie.GenerateRandomKey(64)), //Signing key
	[]byte(securecookie.GenerateRandomKey(32)))

var ErrInvalidPassword = errors.New("Invalid Password")

// Login attempts to login the user given a request.
func Login(r *http.Request) (bool, error) {
	username, password := r.FormValue("username"), r.FormValue("password")
	session, _ := Store.Get(r, "goUntitled")
	u, err := models.GetUserByUsername(username)
	if err != nil && err != models.ErrUsernameTaken {
		return false, err
	}
	//If we've made it here, we should have a valid user stored in u
	//Let's check the password
	err = bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(password))
	if err != nil {
		ctx.Set(r, "user", nil)
		return false, ErrInvalidPassword
	}
	ctx.Set(r, "user", u)
	session.Values["id"] = u.Id
	return true, nil
}

var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

// Register attempts to register the user given a request.
func Register(r *http.Request) (bool, error) {

	Logger.Printf("BODY: %s", r.Body)
	decoder := json.NewDecoder(r.Body)
	reqUser := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := decoder.Decode(&reqUser)
	if err != nil {
		panic("didnt decode")
	}
	username := reqUser.Username
	password := reqUser.Password

	Logger.Printf("%s - %s", username, password)

	//username, password := r.FormValue("username"), r.FormValue("password")
	u, err := models.GetUserByUsername(username)
	// If we have an error which is not simply indicating that no user was found, report it
	if err != sql.ErrNoRows {
		fmt.Println(err)
		return false, err
	}
	//If we've made it here, we should have a valid username given
	//Let's create the password hash
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u.Username = username
	u.Hash = string(h)
	u.ApiKey = GenerateSecureKey()
	if err != nil {
		return false, err
	}
	err = models.Conn.Insert(&u)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GenerateSecureKey() string {
	// Inspired from gorilla/securecookie
	k := make([]byte, 32)
	io.ReadFull(rand.Reader, k)
	return fmt.Sprintf("%x", k)
}

func ChangePassword(r *http.Request) error {
	u := ctx.Get(r, "user").(models.User)
	c, n := r.FormValue("current_password"), r.FormValue("new_password")
	// Check the current password
	err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(c))
	if err != nil {
		return ErrInvalidPassword
	} else {
		// Generate the new hash
		h, err := bcrypt.GenerateFromPassword([]byte(n), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Hash = string(h)
		if err = models.PutUser(&u); err != nil {
			return err
		}
		return nil
	}
}
