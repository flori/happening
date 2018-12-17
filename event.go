package happening

import (
	"time"

	"github.com/lib/pq"
)

type Event struct {
	Id       string         `json:"id" validate:"required" gorm:"type:uuid";"primary_key"`
	Name     string         `json:"name" validate:"required" gorm:"notnull" gorm:"text"`
	Command  pq.StringArray `json:"command,omitempty" gorm:"type:text[]"`
	Output   string         `json:"output,omitempty" gorm:"type:text"`
	Started  time.Time      `json:"started" validate:"required" gorm:"type:timestamptz";"index"`
	Duration time.Duration  `json:"duration" gorm:"type:bigint"`
	Success  bool           `json:"success" gorm:"type:bool";"notnull"`
	ExitCode int            `json:"exitCode" gorm:"type:smallint";"notnull"`
	Hostname string         `json:"hostname" gorm:"type:text"`
	Pid      int            `json:"pid" gorm:"type:int"`
	Store    bool           `json:"store" gorm:"-"`
}
