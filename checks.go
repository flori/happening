package happening

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

func addToChecks(api *API, check *Check) error {
	return api.DB.Create(check).Error
}

func computeHealthStatus(api *API, checks *[]Check) {
	now := time.Now()
	for i, check := range *checks {
		healthBefore := check.Healthy
		time := check.LastPingAt.Add(check.Period)
		healthNow := time.After(now) && check.Success
		(*checks)[i].Healthy = healthNow
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
	var checks []Check
	if err := api.DB.Find(&checks).Error; err != nil {
		log.Panic(err)
		return
	}
	log.Printf("Updating health status of %d checks", len(checks))
	computeHealthStatus(api, &checks)
	for _, check := range checks {
		if err := api.DB.Save(&check).Error; err != nil {
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

func updateCheck(api *API, event *Event) {
	var check Check
	if err := api.DB.Where("name = ?", event.Name).First(&check).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("error: %v", err)
		}
		return
	}
	check.LastPingAt = event.Started
	check.Success = event.Success
	if err := api.DB.Save(&check).Error; err != nil {
		log.Printf("error: %v", err)
	}
}
