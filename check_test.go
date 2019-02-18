package happening

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckInitSuccess(t *testing.T) {
	check := &Check{
		Failures:        0,
		AllowedFailures: 0,
	}
	assert.False(t, check.Success)
	check.Init()
	assert.True(t, check.Success)
}

func TestCheckInitNoSuccess(t *testing.T) {
	check := &Check{
		Failures:        1,
		AllowedFailures: 0,
	}
	assert.False(t, check.Success)
	check.Init()
	assert.False(t, check.Success)
}

func TestCheckInitNoSuccessButAllowed(t *testing.T) {
	check := &Check{
		Failures:        1,
		AllowedFailures: 1,
	}
	assert.False(t, check.Success)
	check.Init()
	assert.True(t, check.Success)
}

func TestCheckInitNoSuccessNotAllowed(t *testing.T) {
	check := &Check{
		Failures:        2,
		AllowedFailures: 1,
	}
	assert.False(t, check.Success)
	check.Init()
	assert.False(t, check.Success)
}
