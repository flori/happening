package happening

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testAPI = &API{
	NOTIFIER: &NullNotifier{},
}

func TestComputeHealthNoSuccess(t *testing.T) {
	checks := []Check{
		Check{
			Healthy:    true,
			Success:    false,
			LastPingAt: time.Now(),
			Period:     time.Second,
		},
	}
	assert.True(t, checks[0].Healthy)
	computeHealthStatus(testAPI, &checks)
	assert.False(t, checks[0].Healthy)
}

func TestComputeHealthSuccess(t *testing.T) {
	checks := []Check{
		Check{
			Healthy:    false,
			Success:    true,
			LastPingAt: time.Now(),
			Period:     time.Second,
		},
	}
	assert.False(t, checks[0].Healthy)
	computeHealthStatus(testAPI, &checks)
	assert.True(t, checks[0].Healthy)
}

func TestComputeHealthTimeout(t *testing.T) {
	oldTime := time.Now().Add(-2 * time.Second)
	checks := []Check{
		Check{
			Healthy:    true,
			Success:    true,
			LastPingAt: oldTime,
			Period:     time.Second,
		},
	}
	assert.True(t, checks[0].Healthy)
	computeHealthStatus(testAPI, &checks)
	assert.False(t, checks[0].Healthy)
}

func TestComputeHealthTimeoutResolved(t *testing.T) {
	checks := []Check{
		Check{
			Healthy:    false,
			Success:    true,
			LastPingAt: time.Now(),
			Period:     2 * time.Second,
		},
	}
	assert.False(t, checks[0].Healthy)
	computeHealthStatus(testAPI, &checks)
	assert.True(t, checks[0].Healthy)
}
