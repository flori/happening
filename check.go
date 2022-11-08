package happening

import (
	"fmt"
	"time"
)

type Check struct {
	Id              *string       `json:"id,omitempty" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Context         string        `json:"context" gorm:"not null;index:,unique,composite:idx_context_name;default:'default'"`
	Name            string        `json:"name" validate:"required" gorm:"type:text;index:,unique,composite:idx_context_name;not null"`
	Period          time.Duration `json:"period" gorm:"type:bigint;not null;default:3600000000000"`
	LastPingAt      time.Time     `json:"last_ping_at" gorm:"type:timestamptz;index;not null;default:now()::timestamptz"`
	LastEventId     *string       `json:"last_event_id,omitempty" gorm:"type:uuid;default:null"`
	Success         bool          `json:"success"`
	Healthy         bool          `json:"healthy" gorm:"type:boolean;not_null;default:true"`
	Disabled        bool          `json:"disabled" gorm:"type:boolean;not_null;default:false"`
	Failures        int           `json:"failures" gorm:"type:int;not null;default:0"`
	AllowedFailures int           `json:"allowed_failures" gorm:"type:int;not null;default:0"`
}

func (check Check) String() string {
	check.Init()
	return fmt.Sprintf(
		`Check:
 - Name: %s
 - Context: %s
 - Id: %s
 - State: %s
 - LastEventId: %s
 - LastPingAt: %s
 - Period: %s
`,
		escapeString(check.Name),
		escapeString(check.Context),
		escapeString(derefString(check.Id)),
		check.State(),
		escapeString(derefString(check.LastEventId)),
		check.LastPingAt,
		check.Period,
	)
}

func (check Check) State() string {
	check.Init()
	result := "healthy"
	if !check.Healthy {
		if check.Success {
			result = "timeout"
		} else {
			result = "failed"
		}
	}
	return result
}

func (check *Check) Init() {
	(*check).Success = check.Failures <= check.AllowedFailures
}
