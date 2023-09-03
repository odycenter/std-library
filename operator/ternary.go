package operator

func IF[T any](b bool, t T, f T) T {
	if b {
		return t
	}
	return f
}
