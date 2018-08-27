package happening

import (
	"encoding/json"
	"log"
)

func addToTimeline(api *API, data []byte) *Event {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		log.Println("error:", err)
		return nil
	}
	if err := api.DB.Insert(&event); err != nil {
		log.Println("error:", err)
		return nil
	}
	return &event
}

func fetchRangeFromTimeline(api *API, q string, filters map[string]string, start int, offset int, limit int) ([]Event, int, error) {
	var events []Event
	var err error

	var total int
	_, err = api.DB.Query(&total, `SELECT COUNT(*) FROM events`)
	if err != nil {
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
	_, err = api.DB.Query(
		&events,
		sql,
		start, start,
		q, "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%",
		filters["id"], "%"+filters["id"]+"%",
		filters["name"], "%"+filters["name"]+"%",
		filters["output"], "%"+filters["output"]+"%",
		filters["hostname"], "%"+filters["hostname"]+"%",
		filters["command"], "%"+filters["command"]+"%",
		filters["ec"], filters["ec"],
		filters["success"], filters["success"],
		offset,
		limit,
	)
	if err != nil {
		log.Println(err)
		return events, 0, err
	}
	return events, total, nil
}
