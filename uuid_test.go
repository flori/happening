package happening

import (
	"regexp"
	"testing"
)

func TestGenerateUUIDv4(t *testing.T) {
	uuid := GenerateUUIDv4()
	matched, err := regexp.MatchString(
		"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", uuid)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if matched {
		t.Logf("regexp ok for %s", uuid)
	} else {
		t.Errorf("regexp not ok for %s", uuid)
	}
}
