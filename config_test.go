package happening

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigCanBeConstructed(t *testing.T) {
	config := NewConfig()
	assert.NotEmpty(t, config.Name)
}
