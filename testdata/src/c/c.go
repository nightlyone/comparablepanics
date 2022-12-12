//go:build go1.20

package c

import "io"

//nolint:unused // collects testdata
func directCompare[T comparable](a, b T) bool {
	return a == b
}

//nolint:unused // collects testdata
func main() {
	var (
		a io.Writer = io.Discard
		b io.Writer = io.Discard
	)
	directCompare(
		a, // want `argument 1/2 \(a\) constrained as comparable may panic at runtime in Go 1.20+`
		b, // want `argument 2/2 \(b\) constrained as comparable may panic at runtime in Go 1.20+`
	)
}
