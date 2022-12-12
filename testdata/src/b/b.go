//go:build go1.20

package a

//nolint:unused // collects testdata
func directCompare[T comparable](a, b T) bool {
	return a == b
}

//nolint:unused // collects testdata
func main() {
	directCompare(
		any("hi"), // want `argument 1/2 constrained as comparable may panic at runtime in Go 1.20+`
		any("hi"), // want `argument 2/2 constrained as comparable may panic at runtime in Go 1.20+`
	)
}
