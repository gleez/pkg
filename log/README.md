
# Mantra log interface + the matching facade for the Go stdlib `log.Logger`

The `mantra/log` library (package `log`) provides an interface for leveled logging so
that any custom logger can be used in exactly the same manner within all parts of an application.
The library further provides a default implementation of the interface based on the Go stdlib
`log.Logger` via a reusable facade.

The library can be installed by one of the following methods:

* using `go get`

	```
	go get github.com/gleez/mantra/log
	```

* via cloning this repository:

	```
	git clone git@github.com:gleez/mantra.git ${GOPATH}/src/github.com/gleez/mantra
	```

## Usage

The package defines the `Log` interface, the `Level` type and the `GoLog` facade for the Go stdlib
`log.Log`, which is used by default.

	import "github.com/gleez/mantra/log

	log.Info("hello %v world", "wonderful")


## License

Copyright (c) 2016 Gleez Technologies, contributors.

Distributed under a MIT style license found in the [LICENSE][license] file.

[go]: https://golang.org
[docs]: https://godoc.org/github.com/gleez/mantra/log
[license]: https://github.com/gleez/mantra/blob/master/LICENSE

Credits
https://github.com/ventu-io/go-log-interface

