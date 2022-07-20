package happening

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
)

const DROP_DATABASE = false

type API struct {
	POSTGRES_URL  string
	DB            *gorm.DB
	NOTIFIER      Notifier
	SERVER_CONFIG ServerConfig
}

func (api *API) PrepareDatabase() {
	log.Println("Opening connection to database server…")
	db, err := gorm.Open("postgres", switchDatabase(api.POSTGRES_URL, "postgres"))
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	databaseName := deriveDatabaseName(api.POSTGRES_URL)

	if DROP_DATABASE {
		if err := db.Exec(fmt.Sprintf(`DROP DATABASE %s`, databaseName)).Error; err != nil {
			log.Println(err)
		}
	}

	dbExists := false
	row := db.Table("pg_database").
		Where("datname = ?", databaseName).
		Select("true").
		Row()
	row.Scan(&dbExists)

	if dbExists {
		log.Printf("Connecting to existent database %s.", databaseName)
	} else {
		log.Printf("Creating nonexistent database %s.", databaseName)
		if err := db.Exec(fmt.Sprintf(`CREATE DATABASE %s`, databaseName)).Error; err != nil {
			log.Panic(err)
		}
	}
	db.Close()

	if db, err = gorm.Open("postgres", api.POSTGRES_URL); err != nil {
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
	gocron.Every(60).Seconds().Do(taskExpireOldEvents, api)
	gocron.Start()
}

type eventsResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
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
		escapeString(event.Name),
		escapeString(event.Started.String()),
		escapeString(event.Duration.String()),
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

func (api *API) GetEventHandler(c echo.Context) error {
	result, event, err := getEvent(api, c.Param("id"))
	switch result {
	case "ok":
		log.Printf(`Get event id=%s`, escapeString(event.Id))
		return c.JSON(
			http.StatusOK,
			eventsResponse{
				Success: true,
				Data:    []Event{*event},
			})
	case "not_found":
		return c.JSON(http.StatusNotFound, eventsResponse{Success: false, Message: err.Error()})
	default:
		return err
	}
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
		for i := 0; i < lenChecks; i++ {
			checks[i].Init()
		}
		log.Printf(`Get %d check`, lenChecks)
		return c.JSON(http.StatusOK, checksResponse{Success: true, Data: checks, Count: &lenChecks, Total: &total})
	} else {
		return err
	}
}

func (api *API) DeleteCheckHandler(c echo.Context) error {
	id := c.Param("id")
	result, err := deleteCheck(api, id)
	switch result {
	case "ok":
		log.Printf(`Delete check id=%s`, escapeString(id))
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
	log.Printf(`Post %s`, escapeString(check.String()))
	if err := c.Validate(check); err != nil {
		return err
	}
	result, err := addToChecks(api, check)
	switch result {
	case "ok":
		log.Printf("Received new check \"%s\"", escapeString(check.Name))
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
	log.Printf(`Patch %s`, escapeString(check.String()))
	if err := c.Validate(check); err != nil {
		return err
	}
	result, err := updateCheck(api, id, check)
	switch result {
	case "ok":
		log.Printf("Received updated check id=\"%s\"", escapeString(id))
		return c.JSON(http.StatusOK, checksResponse{Success: true, Id: id})
	case "not_found":
		return c.JSON(http.StatusNotFound, checksResponse{Success: false, Message: err.Error()})
	default:
		return err
	}
}

func (api *API) ResetCheckHandler(c echo.Context) error {
	id := c.Param("id")
	log.Printf(`Reset check %s`, escapeString(id))
	result, err := resetCheck(api, id)
	switch result {
	case "ok":
		log.Printf("Reset check id=\"%s\"", escapeString(id))
		return c.JSON(http.StatusOK, checksResponse{Success: true, Id: id})
	case "not_found":
		return c.JSON(http.StatusNotFound, checksResponse{Success: false, Message: err.Error()})
	default:
		return err
	}
}

func (api *API) GetCheckHandler(c echo.Context) error {
	result, check, err := getCheck(api, c.Param("id"))
	switch result {
	case "ok":
		log.Printf(`Get check id=%s`, escapeString(*check.Id))
		return c.JSON(
			http.StatusOK,
			checksResponse{
				Success: true,
				Data:    []Check{*check},
			})
	case "not_found":
		return c.JSON(http.StatusNotFound, checksResponse{Success: false, Message: err.Error()})
	default:
		return err
	}
}

func (api *API) GetCheckByNameHandler(c echo.Context) error {
	name := c.Param("name")
	result, check, err := getCheckByName(api, name)
	switch result {
	case "ok":
		log.Printf(`Get check by name "%s", resolved as check id=%s`, escapeString(check.Name), escapeString(*check.Id))
		return c.JSON(
			http.StatusOK,
			checksResponse{
				Success: true,
				Data:    []Check{*check},
			})
	case "not_found":
		return c.JSON(http.StatusNotFound, checksResponse{Success: false, Message: err.Error()})
	default:
		return err
	}
}
