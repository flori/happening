package happening

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventToJSON(t *testing.T) {
	started, _ := time.Parse(time.RFC3339, time.RFC3339)
	event := &Event{
		Id:          "0d7db99e-45ee-4dd8-a637-a30aa90fc3d3",
		Name:        "TestEVent",
		Context:     "default",
		Started:     started,
		Duration:    23,
		Success:     true,
		Hostname:    "localhost",
		Pid:         666,
		Load:        0.5,
		CpuUsage:    0.8,
		MemoryUsage: 6660000,
		Store:       true,
	}
	json := string(EventToJSON(event))
	assert.Equal(t, json, `{"id":"0d7db99e-45ee-4dd8-a637-a30aa90fc3d3","context":"default","name":"TestEVent","started":"0001-01-01T00:00:00Z","duration":23,"success":true,"exitCode":0,"hostname":"localhost","pid":666,"load":0.5,"cpuUsage":0.8,"memoryUsage":6660000,"store":true}`)
}
