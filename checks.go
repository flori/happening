package happening

import (
	"log"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

func wasUniqueViolation(err error) bool {
	if err != nil {
		switch reflect.ValueOf(err).Interface().(type) {
		case *pq.Error:
			pqError := err.(*pq.Error)
			return pqError.Code.Name() == "unique_violation"
		}
	}
	return false
}

func addToChecks(api *API, check *Check) (string, error) {
	err := api.DB.Create(check).Error
	if err == nil {
		return "ok", err
	}
	if wasUniqueViolation(err) {
		return "conflict", err
	}
	return "error", err
}

func updateCheck(api *API, id string, check *Check) (string, error) {
	checkInstance := api.DB.Model(&Check{Id: &id})
	checkInstance = checkInstance.Update(
		"disabled",
		check.Disabled,
	).Update(
		"period",
		check.Period,
	).Update(
		"allowed_failures",
		check.AllowedFailures,
	)

	if check.Healthy && check.Failures == 0 && check.Success {
		checkInstance = checkInstance.Update(
			"healthy",
			check.Healthy,
		).Update(
			"failures",
			0,
		).Update(
			"success",
			true,
		).Update(
			"last_ping_at",
			time.Now(),
		)
	}

	if err := checkInstance.Error; err != nil {
		return "not_found", err
	}
	return "ok", nil
}

func resetCheck(api *API, id string) (string, error) {
	checkInstance := api.DB.Model(&Check{Id: &id})
	checkInstance = checkInstance.Update(
		"healthy",
		true,
	).Update(
		"success",
		true,
	).Update(
		"last_ping_at",
		time.Now(),
	).Update(
		"failures",
		0,
	)

	if err := checkInstance.Error; err != nil {
		return "not_found", err
	}
	return "ok", nil
}

func deleteCheck(api *API, id string) (string, error) {
	var check Check
	if err := api.DB.Where("id = ?", id).First(&check).Error; err != nil {
		return "not_found", err
	}
	if err := api.DB.Delete(&check).Error; err != nil {
		return "error", err
	}
	return "ok", nil
}

func getCheck(api *API, id string) (string, *Check, error) {
	var check Check
	if err := api.DB.Where("id = ?", id).First(&check).Error; err != nil {
		return "not_found", nil, err
	}
	check.Init()
	return "ok", &check, nil
}

func getCheckByNameInContext(api *API, name string, context string) (string, *Check, error) {
	var check Check
	if err := api.DB.Where("name = ?", name).
		Where("context = ?", context).First(&check).Error; err != nil {
		return "not_found", nil, err
	}
	check.Init()
	return "ok", &check, nil
}

func computeHealthStatus(api *API, checks *[]Check) {
	now := time.Now()
	for i, check := range *checks {
		healthBefore := check.Healthy
		time := check.LastPingAt.Add(check.Period)
		check.Init()
		healthNow := time.After(now) && check.Success
		(*checks)[i].Healthy = healthNow
		if check.Disabled {
			continue
		}
		if healthBefore && !healthNow {
			log.Printf("Alert: %s", (*checks)[i])
			api.NOTIFIER.Alert((*checks)[i])
		}
		if !healthBefore && healthNow {
			log.Printf("Resolve: %s", (*checks)[i])
			api.NOTIFIER.Resolve((*checks)[i])
		}
	}
}

func taskUpdateHealthStatus(api *API) {
	tx := api.DB.Begin()
	tx.Exec("LOCK TABLE checks")
	var checks []Check
	if err := tx.Find(&checks).Error; err != nil {
		log.Panic(err)
		return
	}
	defer tx.Commit()
	computeHealthStatus(api, &checks)
	for _, check := range checks {
		if err := tx.Save(&check).Error; err != nil {
			log.Printf("error: %v", err)
		}
	}
}

func fetchRangeFromChecks(api *API, p parameters) ([]Check, int, error) {
	var checks []Check

	var total int
	if err := api.DB.Model(&Check{}).Count(&total).Error; err != nil {
		return checks, 0, err
	}
	if err := api.DB.Model(&Check{}).Order("name ASC").Offset(p.Offset).Limit(p.Limit).
		Scan(&checks).Error; err != nil {
		log.Println(err)
		return checks, 0, err
	}
	return checks, total, nil
}

func updateCheckOnEvent(api *API, event *Event) {
	var check Check
	if err := api.DB.Where("name = ?", event.Name).
		Where("context = ?", event.Context).
		First(&check).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("error: %v", err)
		}
		return
	}
	check.LastPingAt = event.Started
	if event.Store {
		check.LastEventId = &event.Id
	} else {
		check.LastEventId = nil
	}
	if event.Success {
		check.Failures = 0
	} else {
		check.Failures++
	}
	if err := api.DB.Save(&check).Error; err != nil {
		log.Printf("error: %v", err)
	}
}
