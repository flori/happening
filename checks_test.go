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
			Name:       "test",
			Context:    "default",
			Healthy:    true,
			Failures:   1,
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
			Name:       "test",
			Context:    "default",
			Healthy:    false,
			Failures:   0,
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
			Name:       "test",
			Context:    "default",
			Healthy:    true,
			Failures:   0,
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
			Name:       "test",
			Context:    "default",
			Healthy:    false,
			Failures:   0,
			LastPingAt: time.Now(),
			Period:     2 * time.Second,
		},
	}
	assert.False(t, checks[0].Healthy)
	computeHealthStatus(testAPI, &checks)
	assert.True(t, checks[0].Healthy)
}

func TestComputeHealthSuccessAfterFailure(t *testing.T) {
	checks := []Check{
		Check{
			Name:            "test",
			Context:         "default",
			Healthy:         true,
			Failures:        1,
			LastPingAt:      time.Now(),
			Period:          time.Second,
			AllowedFailures: 1,
		},
	}
	assert.True(t, checks[0].Healthy)
	computeHealthStatus(testAPI, &checks)
	assert.True(t, checks[0].Healthy)
}

func TestComputeHealthNoSuccessAfterFailure(t *testing.T) {
	checks := []Check{
		Check{
			Name:            "test",
			Context:         "default",
			Healthy:         true,
			Failures:        2,
			LastPingAt:      time.Now(),
			Period:          time.Second,
			AllowedFailures: 1,
		},
	}
	assert.True(t, checks[0].Healthy)
	computeHealthStatus(testAPI, &checks)
	assert.False(t, checks[0].Healthy)
}
