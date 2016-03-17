
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

The package defines the `Logger` interface, the `Level` type and the `GoLog` facade for the Go stdlib
`log.Logger`, which is used by default. The default logger can be obtained via a "factory" function
`GetLogger`. The only accepted argument to the latter, `name`, is not used at the moment, but is
defined for future API compatibility:

	import "github.com/gleez/mantra/log

	logger := log.GetLogger("root")
	logger.Info("hello %v world", "wonderful")

By calling `SetLogger` and giving a custom implementation of the `Logger` interface, one
can adopt any other logging framework (here using the implementation from the library itself):

	import "github.com/gleez/mantra/log

	log.SetLogger(log.NewGoLog(golog.New(tw, "", golog.Ldate|golog.Ltime|golog.LUTC), log.INFO, "APP", false))

	logger := log.GetLogger("root")
	logger.Info("hello %v world", "wonderful")

The `Logger` interface defines the following self explaining methods:

	type Logger interface {
        // SetLevel sets the logging level of the logger. Messages with the level lower than currently set
        // are ignored. Messages with the level FATAL or PANIC are never ignored.
        SetLevel(level Level)

        // GetLevel retrieves the currently set logger level.
        GetLevel() Level

        // Log logs a formatted message with a given level. If level is below the one currently set
        // the message will be ignored.
        Log(level Level, message string, args ...interface{})

        // Trace logs a formatted message with the TRACE level.
        Trace(message string, args ...interface{})

        // Debug logs a formatted message with the DEBUG level.
        Debug(message string, args ...interface{})

        // Info logs a formatted message with the INFO level.
        Info(message string, args ...interface{})

		// Notice logs a formatted message with the NOTICE level.
		Notice(message string, args ...interface{})
	
        // Warn logs a formatted message with the WARN level.
        Warn(message string, args ...interface{})

        // Error logs a formatted message with the ERROR level.
        Error(message string, args ...interface{})

        // Panic logs a formatted message with the PANIC level and calls panic.
        Panic(message string, args ...interface{})

		// Alert logs a formatted message with the ALERT level.
		Alert(message string, args ...interface{})
	
        // Fatal logs a formatted message with the FATAL level and exits the application with an error.
        Fatal(message string, args ...interface{})
    }

## License

Copyright (c) 2016 Gleez Technologies, contributors.

Distributed under a MIT style license found in the [LICENSE][license] file.

[go]: https://golang.org
[docs]: https://godoc.org/github.com/gleez/mantra/log
[license]: https://github.com/gleez/mantra/blob/master/LICENSE

Credits
https://github.com/ventu-io/go-log-interface

