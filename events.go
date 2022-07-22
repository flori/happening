package happening

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func addToEvents(api *API, event *Event) error {
	return api.DB.Create(event).Error
}

func fetchRangeFromEvents(api *API, p parameters) ([]Event, int, error) {
	var events []Event

	var total int
	if err := api.DB.Model(&Event{}).Count(&total).Error; err != nil {
		log.Panic(err)
		return events, 0, err
	}
	sql := `
		SELECT *
		FROM events WHERE
			(? = 0 OR started > TO_TIMESTAMP(?)) AND
			(? = '' OR (name LIKE ?) OR (output LIKE ?) OR (hostname LIKE ?) OR (CAST(command AS text) LIKE ?)) AND
			(? = '' OR CAST(id AS TEXT) LIKE ?) AND
			(? = '' OR name LIKE ?) AND
			(? = '' OR output LIKE ?) AND
			(? = '' OR hostname LIKE ?) AND
			(? = '' OR CAST(command AS text) LIKE ?) AND
			(? = '' OR CAST(exit_code AS text) = ?) AND
			(? = '' OR success = (? = 'true')) AND
			(? = '' OR context = ?)
		ORDER BY started DESC
		OFFSET ?
		LIMIT ?
	`
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
		p.Filters["context"], p.Filters["context"],
		p.Offset,
		p.Limit,
	).Scan(&events)
	if err := api.DB.Error; err != nil {
		log.Println(err)
		return []Event{}, 0, err
	}
	for i := 0; i < len(events); i++ {
		event := &events[i]
		if event.Output == "" {
			continue
		}
		lenOutput := strings.Count(event.Output, "\n")
		if lenOutput == 0 {
			lenOutput = 1
		}
		event.Output = fmt.Sprintf("%d lines", lenOutput)
	}
	return events, total, nil
}

func getEvent(api *API, id string) (string, *Event, error) {
	var event Event
	if err := api.DB.Where("id = ?", id).First(&event).Error; err != nil {
		return "not_found", nil, err
	}
	return "ok", &event, nil
}

func taskExpireOldEvents(api *API) {
	seconds, err := strconv.Atoi(api.SERVER_CONFIG.RETENTION_DAYS)
	if err != nil || seconds <= 0 {
		return
	}
	retentionPeriod := time.Duration(seconds) * 24 * time.Hour
	expireTimestamp := time.Now().Add(-retentionPeriod)
	log.Printf("Now expiring events before %s", expireTimestamp.Format(time.RFC3339))
	err = api.DB.Where("started < ?", expireTimestamp).Delete(&Event{}).Error
	if err != nil {
		log.Panic(err)
	}
}
