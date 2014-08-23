package models

import (
	"errors"
	"log"
	"os"

	"github.com/coopernurse/gorp"
	"github.com/daragao/goUntitled/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var Conn *gorp.DbMap
var db gorm.DB
var err error
var ErrUsernameTaken = errors.New("username already taken")
var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

const (
	CAMPAIGN_IN_PROGRESS string = "In progress"
	CAMPAIGN_QUEUED      string = "Queued"
	CAMPAIGN_COMPLETE    string = "Completed"
	EVENT_SENT           string = "Email Sent"
	EVENT_OPENED         string = "Email Opened"
	EVENT_CLICKED        string = "Clicked Link"
	STATUS_SUCCESS       string = "Success"
	STATUS_UNKNOWN       string = "Unknown"
	ERROR                string = "Error"
)

// Flash is used to hold flash information for use in templates.
type Flash struct {
	Type    string
	Message string
}

type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func Setup() error {
	Logger.Printf("Setup for models.")
	//db, err = gorm.Open("sqlite3", config.Conf.DBPath)
	databaseURL := config.Conf.Database.User + ":" +
		config.Conf.Database.Password + "@/" +
		config.Conf.Database.Name + "?charset=utf8&parseTime=True"
	db, err = gorm.Open("mysql", databaseURL)
	db.LogMode(false)
	db.SetLogger(Logger)
	if err != nil {
		Logger.Println(err)
		return err
	}
	//If the file already exists, delete it and recreate it
	//_, err = os.Stat(config.Conf.DBPath)
	if err == nil {
		Logger.Printf("Database not found... creating db at %s\n",
			config.Conf.Database.Host)
		db.DropTableIfExists(User{})
		db.CreateTable(User{})
		/*db.CreateTable(Target{})
		db.CreateTable(Result{})
		db.CreateTable(Group{})
		db.CreateTable(GroupTarget{})
		db.CreateTable(Template{})
		db.CreateTable(Attachment{})
		db.CreateTable(SMTP{})
		db.CreateTable(Event{})
		db.CreateTable(Campaign{})*/
		//Create the default user
		init_user := User{
			Username: "admin",
			Hash:     "$2a$10$IYkPp0.QsM81lYYPrQx6W.U6oQGw7wMpozrKhKAHUBVL4mkm/EvAS", //gophish
			ApiKey:   "12345678901234567890123456789012",
		}
		err = db.Save(&init_user).Error
		if err != nil {
			Logger.Println(err)
		}
	}
	return nil
}
