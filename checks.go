package happening

import (
	"encoding/json"
	"log"
	"time"
)

func addToChecks(api *API, data []byte) *Check {
	var check Check
	log.Println(string(data))
	err := json.Unmarshal(data, &check)
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}
	if err := api.DB.Create(&check).Error; err != nil {
		log.Printf("error: %v", err)
		return nil
	}
	return &check
}

func computeHealthStatus(api *API, checks *[]Check) {
	now := time.Now()
	for i, check := range *checks {
		healthBefore := check.Healthy
		time := check.LastPingAt.Add(check.Period)
		healthNow := time.After(now)
		(*checks)[i].Healthy = healthNow
		if healthBefore && !healthNow {
			log.Println((*checks)[i])
			api.NOTIFIER.Alert((*checks)[i])
		}
	}
}

func taskUpdateHealthStatus(api *API) {
	var checks []Check
	if err := api.DB.Find(&checks).Error; err != nil {
		log.Printf("foo error: %v", err)
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
	if !event.Success {
		return
	}
	var check Check
	if err := api.DB.Where("name = ?", event.Name).First(&check).Error; err != nil {
		log.Printf("error: %v", err)
		return
	}
	check.LastPingAt = event.Started
	var checks []Check
	checks = append(checks, check)
	computeHealthStatus(api, &checks)
	if err := api.DB.Save(&check).Error; err != nil {
		log.Printf("error: %v", err)
	}
}
