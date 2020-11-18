# apexlogutils

Go packages with utilities around github.com/apex/log.

## What does it do?

This code was extracted from multiple production applications for integration with *HTTP request logging*, *DB driver logging with pgx* and a custom text log handler that
outputs a component field to quickly identify the source of a log message in development.

The HTTP logger is based on https://github.com/apex/httplog but is a more complete `http.ResponseWriter` wrapper (see https://github.com/gorilla/handlers) and includes an
option to not log requests for paths with a given prefix (e.g. health checks). It also uses a log instance from `context.Context` if set, so it's a nice combination with
`middleware.RequestID` to have a unique id for each request to correlate log messages.
