package happening

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUUIDv4(t *testing.T) {
	uuid := GenerateUUIDv4()
	assert.Regexp(t, "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", uuid)
}
