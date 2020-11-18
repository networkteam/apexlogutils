package pgx

import (
	"github.com/jackc/pgx/v4"
	"github.com/networkteam/apexlogutils"
)

// ToPgxLogLevel gets the pgx log level from Verbosity
func ToPgxLogLevel(v apexlogutils.Verbosity) pgx.LogLevel {
	// Pgx log levels are somehow borked, queries and connections are only logged with trace level
	if v >= apexlogutils.VerbosityDebug {
		return pgx.LogLevelTrace
	}

	switch v {
	case apexlogutils.VerbosityInfo:
		return pgx.LogLevelDebug
	case apexlogutils.VerbosityWarn:
		return pgx.LogLevelInfo
	case apexlogutils.VerbosityError:
		return pgx.LogLevelError
	}

	return pgx.LogLevelNone
}
