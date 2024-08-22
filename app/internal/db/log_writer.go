package internal_db

type DoNothingWriter struct {
}

func (sw *DoNothingWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
