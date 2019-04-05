package happening

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testAPI = &API{
	NOTIFIER: &NullNotifier{},
}

func TestComputeHealthStatus(t *testing.T) {
	checks := new([]Check)
	computeHealthStatus(testAPI, checks)
	assert.True(t, true)
}
