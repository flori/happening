package happening

import (
	"fmt"
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
	log.Println("Opening connection to database server…")
	postgresURL := fmt.Sprintf(api.POSTGRES_URL, "postgres")
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
	event := new(Event)
	if err := c.Bind(event); err != nil {
		return err
	}
	if err := c.Validate(event); err != nil {
		return err
	}
	if event.Store {
		if err := addToEvents(api, event); err != nil {
			return err
		}
	}
	updateCheckOnEvent(api, event)
	log.Printf(
		"Received event \"%s\", started %v, lasted %v\n",
		event.Name,
		event.Started,
		event.Duration,
	)
	return c.JSON(http.StatusOK, eventsResponse{Success: true})
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
	check := new(Check)
	if err := c.Bind(check); err != nil {
		return err
	}
	log.Printf(`Post %s`, *check)
	if err := c.Validate(check); err != nil {
		return err
	}
	result, err := addToChecks(api, check)
	switch result {
	case "ok":
		log.Printf("Received new check \"%s\"", check.Name)
		return c.JSON(http.StatusOK, checksResponse{Success: true})
	case "conflict":
		return c.JSON(http.StatusConflict, checksResponse{Success: false, Message: err.Error()})
	default:
		return err
	}
}

func (api *API) PatchCheckHandler(c echo.Context) error {
	id := c.Param("id")
	check := new(Check)
	if err := c.Bind(check); err != nil {
		return err
	}
	log.Printf(`Patch %s`, *check)
	if err := c.Validate(check); err != nil {
		return err
	}
	if err := c.Validate(check); err != nil {
		return err
	}
	result, err := updateCheck(api, id, check)
	switch result {
	case "ok":
		log.Printf("Received updated check id=\"%s\"", id)
		return c.JSON(http.StatusOK, checksResponse{Success: true, Id: id})
	case "not_found":
		return c.JSON(http.StatusNotFound, checksResponse{Success: false, Message: err.Error()})
	default:
		return err
	}
}
