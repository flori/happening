package happening

import "strings"

func derefString(s *string) string {
	result := "<nil>"
	if s != nil {
		result = *s
	}
	return result
}

func escapeString(s string) string {
	escaped := strings.Replace(s, "\n", "\\n", -1)
	escaped = strings.Replace(escaped, "\r", "\\r", -1)
	return escaped
}
