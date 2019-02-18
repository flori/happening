package happening

import (
	"fmt"
	"time"
)

type Check struct {
	Id              *string       `json:"id,omitempty" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name            string        `json:"name" validate:"required" gorm:"type:text;unique;not null"`
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
	return fmt.Sprintf(
		`Check:
 - Name: %s
 - Id: %s
 - State: %s
 - LastEventId: %s
 - LastPingAt: %s
 - Period: %s
`,
		check.Name,
		derefString(check.Id),
		check.State(),
		derefString(check.LastEventId),
		check.LastPingAt,
		check.Period,
	)
}

func (check Check) State() string {
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
