package happening

import (
	"io/ioutil"
	"log"
	"net/url"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
)

const MaxInt = int(^uint(0) >> 1)

type API struct {
	POSTGRES_URL string
	DB           *pg.DB
}

type response struct {
	Success bool    `json:"success"`
	Count   int     `json:"count"`
	Total   int     `json:"total"`
	Data    []Event `json:"data,omitempty"`
}

func (api *API) PrepareDatabase() {
	dbOptions, err := pg.ParseURL(api.POSTGRES_URL)
	if err != nil {
		log.Panic(err)
	}
	db := pg.Connect(dbOptions)
	defer db.Close()
	dbExists := false
	_, err = db.Query(&dbExists, `SELECT true FROM pg_database WHERE datname = 'happening'`)
	if err != nil {
		log.Panic(err)
	}
	if dbExists {
		db.Close()
		dbOptions.Database = "happening"
		api.DB = pg.Connect(dbOptions)
	} else {
		_, err = db.Exec(`CREATE DATABASE happening`)
		if err != nil {
			log.Panic(err)
		}
		db.Close()
		dbOptions.Database = "happening"
		db = pg.Connect(dbOptions)
		for _, model := range []interface{}{&Event{}} {
			err := db.CreateTable(model, &orm.CreateTableOptions{})
			if err != nil {
				log.Panic(err)
			}
		}
		_, err = db.Exec(`CREATE INDEX started_idx ON events (started DESC)`)
		if err != nil {
			log.Panic(err)
		}
		api.DB = db
	}
}

// PostEventHandler handles the posting of new events
func (api *API) PostEventHandler(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Println("error:", err)
	} else {
		if event := addToTimeline(api, data); event != nil {
			log.Printf("info: Received event \"%s\", started %v, lasted %v\n",
				event.Name, event.Started, event.Duration)
			return c.JSON(200, response{Success: true})
		}
	}
	return c.JSON(500, response{Success: false})
}

func parseFilters(url *url.URL) map[string]string {
	filters := make(map[string]string)
	for param, values := range url.Query() {
		if len(param) > 2 && param[0:2] == "f:" {
			filters[param[2:]] = values[len(values)-1]
		}
	}
	return filters
}

func parseParameters(url *url.URL) (string, map[string]string, int, int, int) {
	query := url.Query().Get("q")
	filters := parseFilters(url)
	var start int
	var offset int
	var limit int
	var err error
	if o := url.Query().Get("o"); o == "" {
		offset = 0
	} else {
		if offset, err = strconv.Atoi(o); err != nil {
			offset = 0
		}
	}
	if l := url.Query().Get("l"); l == "" {
		limit = 50
	} else {
		if l == "*" {
			limit = MaxInt
		} else if limit, err = strconv.Atoi(l); err != nil {
			limit = 50
		}
	}
	if s := url.Query().Get("s"); s != "" {
		if start, err = strconv.Atoi(s); err != nil {
			start = 0
		}
	}
	return query, filters, start, offset, limit
}

// GetEventsHandler handles listing of events
func (api *API) GetEventsHandler(c echo.Context) error {
	query, filters, start, offset, limit := parseParameters(c.Request().URL)
	events, total, err := fetchRangeFromTimeline(api, query, filters, start, offset, limit)
	if err == nil {
		return c.JSON(200, response{Success: true, Data: events, Count: len(events), Total: total})
	}
	return c.JSON(500, response{Success: false})
}
