package happening

func derefString(s *string) string {
	result := "<nil>"
	if s != nil {
		result = *s
	}
	return result
}
