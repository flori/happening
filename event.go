package happening

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Event struct {
	Id          string         `json:"id" validate:"required,uuid" gorm:"type:uuid;primary_key"`
	Context     string         `json:"context" gorm:"not null;index;default:'default'"`
	Name        string         `json:"name" validate:"required,printascii" gorm:"not null;text"`
	Command     pq.StringArray `json:"command,omitempty" gorm:"type:text[]"`
	Output      string         `json:"output,omitempty" gorm:"type:text"`
	Started     time.Time      `json:"started" validate:"required" gorm:"type:timestamptz;index"`
	Duration    time.Duration  `json:"duration" gorm:"type:bigint"`
	Success     bool           `json:"success" gorm:"type:bool;not null"`
	ExitCode    int            `json:"exitCode" gorm:"type:smallint;not null"`
	Signal      string         `json:"signal,omitempty" gorm:"text"`
	Hostname    string         `json:"hostname" gorm:"type:text"`
	User        string         `json:"user", gorm:"type:text"`
	Pid         int            `json:"pid" gorm:"type:int"`
	Load        float32        `json:"load" gorm:"type:real"`
	CpuUsage    float64        `json:"cpuUsage" gorm:"type:real"`
	MemoryUsage float64        `json:"memoryUsage" gorm:"type:real"`
	Store       bool           `json:"store" gorm:"-"`
}

func (event Event) Result() string {
	result := "failure"
	if event.Success {
		result = "success"
	}
	return result
}

func (event Event) String() string {
	return fmt.Sprintf(
		`Event:
 - Name: %s
 - Context: %s
 - Id: %s
 - Result: %s
 - Started: %s
 - Duration: %s
 - Hostname: %s
`,
		escapeString(event.Name),
		escapeString(event.Context),
		escapeString(event.Id),
		event.Result(),
		event.Started,
		event.Duration,
		event.Hostname,
	)
}
