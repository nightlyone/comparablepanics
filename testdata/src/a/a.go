package a

//nolint:unused // collects testdata
func directCompare[T comparable](a, b T) bool {
	return a == b
}

//nolint:unused // collects testdata
func main() {
	directCompare(int(1), int(1))
}
