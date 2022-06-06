package happening

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDerefString(t *testing.T) {
	var s *string
	assert.Equal(t, "<nil>", derefString(s))
	s1 := "foo"
	s = &s1
	assert.Equal(t, "foo", derefString(s))
}

func TestEscapeString(t *testing.T) {
	assert.Equal(t, "foo\\rbar\\nquux", escapeString("foo\rbar\nquux"))
}
