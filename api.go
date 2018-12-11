package happening

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
)

const DEBUG = false

const MaxInt = int(^uint(0) >> 1)

type API struct {
	POSTGRES_URL  string
	DATABASE_NAME string
	DB            *gorm.DB
	NOTIFIER      Notifier
	SERVER_CONFIG ServerConfig
}

func (api *API) PrepareDatabase() {
	postgresURL := fmt.Sprintf(api.POSTGRES_URL, "postgres")
	log.Printf("Opening connection to database server %s…", postgresURL)
	db, err := gorm.Open("postgres", postgresURL)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	if DEBUG {
		if err := db.Exec(fmt.Sprintf(`DROP DATABASE %s`, api.DATABASE_NAME)).Error; err != nil {
			log.Println(err)
		}
	}

	dbExists := false
	row := db.Table("pg_database").
		Where("datname = ?", api.DATABASE_NAME).
		Select("true").
		Row()
	row.Scan(&dbExists)

	if dbExists {
		log.Printf("Connecting to existent database %s.", api.DATABASE_NAME)
	} else {
		log.Printf("Creating nonexistent database %s.", api.DATABASE_NAME)
		if err := db.Exec(fmt.Sprintf(`CREATE DATABASE %s`, api.DATABASE_NAME)).Error; err != nil {
			log.Panic(err)
		}
	}
	db.Close()

	if db, err = gorm.Open("postgres", fmt.Sprintf(api.POSTGRES_URL, api.DATABASE_NAME)); err != nil {
		log.Panic(err)
	}
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public`).Error; err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&Event{}, &Check{})
	api.DB = db
}

func (api *API) SetupCronJobs() {
	gocron.Every(5).Seconds().Do(taskUpdateHealthStatus, api)
	gocron.Start()
}

type eventsResponse struct {
	Success bool    `json:"success"`
	Count   int     `json:"count"`
	Total   int     `json:"total"`
	Data    []Event `json:"data,omitempty"`
}

// PostEventHandler handles the posting of new events
func (api *API) PostEventHandler(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	} else {
		if event := addToEvents(api, data); event != nil {
			log.Printf("info: Received event \"%s\", started %v, lasted %v\n",
				event.Name, event.Started, event.Duration)
			updateCheck(api, event)
			return c.JSON(http.StatusOK, eventsResponse{Success: true})
		}
	}
	return c.JSON(http.StatusInternalServerError, eventsResponse{Success: false})
}

// GetEventsHandler handles listing of events
func (api *API) GetEventsHandler(c echo.Context) error {
	p := parseParameters(c)
	events, total, err := fetchRangeFromEvents(api, p)
	if err == nil {
		return c.JSON(http.StatusOK, eventsResponse{Success: true, Data: events, Count: len(events), Total: total})
	}
	return c.JSON(http.StatusInternalServerError, eventsResponse{Success: false})
}

type checksResponse struct {
	Success bool    `json:"success"`
	Id      string  `json:"id,omitempty"`
	Message string  `json:"message,omitempty"`
	Count   *int    `json:"count,omitempty"`
	Total   *int    `json:"total,omitempty"`
	Data    []Check `json:"data,omitempty"`
}

// GetChecksHandler handles listing of checks
func (api *API) GetChecksHandler(c echo.Context) error {
	p := parseParameters(c)
	checks, total, err := fetchRangeFromChecks(api, p)
	if err == nil {
		lenChecks := len(checks)
		return c.JSON(http.StatusOK, checksResponse{Success: true, Data: checks, Count: &lenChecks, Total: &total})
	} else {
		return err
	}
}

func (api *API) DeleteCheckHandler(c echo.Context) error {
	result, err := deleteCheck(api, c.Param("id"))
	switch result {
	case "ok":
		return c.JSON(http.StatusOK, checksResponse{Success: true})
	case "not_found":
		return c.JSON(http.StatusNotFound, checksResponse{Success: false, Message: err.Error()})
	default:
		return err
	}
}

func (api *API) PostCheckHandler(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	} else {
		check, err := addToChecks(api, data)
		if err != nil {
			log.Printf("info: Received new check \"%s\"", check.Name)
			return c.JSON(http.StatusOK, checksResponse{Success: true, Id: check.Id})
		} else {
			log.Printf("foo %v", err)
			return err
		}
	}
}
