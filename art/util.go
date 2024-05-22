package art

func makePointer[T any](thing T) *T {
	return &thing
}
