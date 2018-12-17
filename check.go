package happening

import (
	"fmt"
	"time"
)

type Check struct {
	Id         string        `json:"id,omitempty" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name       string        `json:"name" validate:"required" gorm:"type:text;unique;not null"`
	Period     time.Duration `json:"period" gorm:"type:bigint;not null;default:3600000000000"`
	LastPingAt time.Time     `json:"last_ping_at" gorm:"type:timestamptz;index;not null;default:now()::timestamptz"`
	Healthy    bool          `json:"healthy" gorm:"type:boolean;not_null;default:true"`
}

func (check Check) String() string {
	if check.Healthy {
		return fmt.Sprintf(`Check "%v" (%v) is healthy, was last pinged at %v less than %v ago.`,
			check.Name,
			check.Id,
			check.LastPingAt,
			check.Period,
		)
	} else {
		return fmt.Sprintf(`Check "%v" (%v) is unhealthy, was last pinged at %v longer than %v ago.`,
			check.Name,
			check.Id,
			check.LastPingAt,
			check.Period,
		)
	}
}
