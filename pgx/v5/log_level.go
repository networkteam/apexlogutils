package v5

import (
	"github.com/jackc/pgx/v5/tracelog"

	"github.com/networkteam/apexlogutils"
)

// ToPgxLogLevel gets the pgx log level from Verbosity
func ToPgxLogLevel(v apexlogutils.Verbosity) tracelog.LogLevel {
	// Pgx log levels are somehow borked, queries and connections are only logged with trace level
	if v >= apexlogutils.VerbosityDebug {
		return tracelog.LogLevelTrace
	}

	switch v {
	case apexlogutils.VerbosityInfo:
		return tracelog.LogLevelDebug
	case apexlogutils.VerbosityWarn:
		return tracelog.LogLevelInfo
	case apexlogutils.VerbosityError:
		return tracelog.LogLevelError
	}

	return tracelog.LogLevelNone
}
