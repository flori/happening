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
	Success    bool          `json:"success" gorm:"type:boolean;not_null;default:true"`
}

func (check Check) String() string {
	healtyhString := "unhealthy"
	if check.Healthy {
		healtyhString = "healthy"
	}
	successString := "failure"
	if check.Success {
		successString = "success"
	}
	return fmt.Sprintf(`Check "%v" (%v) is %s, was last pinged at %v (period=%s) and was a "%s".`,
		check.Name,
		check.Id,
		healtyhString,
		check.LastPingAt,
		check.Period,
		successString,
	)
}
