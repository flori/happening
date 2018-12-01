package happening

import (
	"encoding/json"
	"log"
)

func addToEvents(api *API, data []byte) *Event {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}
	if event.Store {
		if err := api.DB.Create(&event).Error; err != nil {
			log.Printf("error: %v", err)
			return nil
		}
	}
	return &event
}

func fetchRangeFromEvents(api *API, p parameters) ([]Event, int, error) {
	var events []Event

	var total int
	if err := api.DB.Model(&Event{}).Count(&total).Error; err != nil {
		log.Println(err)
		return events, 0, err
	}
	sql :=
		`SELECT * FROM events WHERE
			(? = 0 OR started > TO_TIMESTAMP(?)) AND
			(? = '' OR (name LIKE ?) OR (output LIKE ?) OR (hostname LIKE ?) OR (CAST(command AS text) LIKE ?)) AND
			(? = '' OR CAST(id AS TEXT) LIKE ?) AND
			(? = '' OR name LIKE ?) AND
			(? = '' OR output LIKE ?) AND
			(? = '' OR hostname LIKE ?) AND
			(? = '' OR CAST(command AS text) LIKE ?) AND
			(? = '' OR CAST(exit_code AS text) = ?) AND
			(? = '' OR success = (? = 'true'))
		ORDER BY started DESC
		OFFSET ?
		LIMIT ?`
	api.DB.Raw(
		sql,
		p.Start, p.Start,
		p.Query, "%"+p.Query+"%", "%"+p.Query+"%", "%"+p.Query+"%", "%"+p.Query+"%",
		p.Filters["id"], "%"+p.Filters["id"]+"%",
		p.Filters["name"], "%"+p.Filters["name"]+"%",
		p.Filters["output"], "%"+p.Filters["output"]+"%",
		p.Filters["hostname"], "%"+p.Filters["hostname"]+"%",
		p.Filters["command"], "%"+p.Filters["command"]+"%",
		p.Filters["ec"], p.Filters["ec"],
		p.Filters["success"], p.Filters["success"],
		p.Offset,
		p.Limit,
	).Scan(&events)
	if err := api.DB.Error; err != nil {
		log.Println(err)
		return events, 0, err
	}
	return events, total, nil
}
