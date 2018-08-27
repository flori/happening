package happening

import (
	"time"
)

type Event struct {
	Id       string        `json:"id" sql:"type:uuid,pk"`
	Name     string        `json:"name"`
	Command  []string      `json:"command,omitempty"`
	Output   string        `json:"output,omitempty"`
	Started  time.Time     `json:"started" sql:"type:timestamptz"`
	Duration time.Duration `json:"duration" sql:"type:bigint"`
	Success  bool          `json:"success" sql:",notnull"`
	ExitCode int           `json:"exitCode" sql:"type:smallint,notnull"`
	Hostname string        `json:"hostname"`
	Pid      int           `json:"pid" sql:"type:bigint`
}
