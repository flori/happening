package happening

type null struct {
}

var NullWriter = null{}

func (null) Write(p []byte) (n int, err error) {
	return len(p), nil
}
