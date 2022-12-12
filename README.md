# Background
[![Go Reference](https://pkg.go.dev/badge/github.com/nightlyone/comparablepanics.svg)](https://pkg.go.dev/github.com/nightlyone/comparablepanics)

Go 1.20+ will change the semantics of the comparable constraint will include
interfaces. Interfaces in Go have always behaved that that way and the
comparable constraint was the outlier.

So far that makes sense, but some people already wrote type safe
generic code that worked with the comparable constraint and may want to know
where their code will now have runtime panics.

This analyzer will mark places where Go 1.20+ comparable will panic to document
the effect of your upgrade and help Go programmers to workaround it.

Another use case is never allowing this new feature in your own code.

## TODO
* [x] Tests (couldn't publish mine for copyright reasons)
* [x] CI/CD with github actions
* [ ] Announce it in golang-nuts
* [ ] Be smarter and follow type inference further
* [ ] Ask for integration into golangci-lint

## Contribution
It is still work in progress and contributions are welcome. Especially test
cases!
