package apexlogutils

import "github.com/apex/log"

// Verbosity is a verbosity level that can be mapped to log levels.
// Higher numbers mean more verbose logging.
type Verbosity int

const (
	// VerbosityFatal is least verbose
	VerbosityFatal Verbosity = 0
	// VerbosityError only logs errors
	VerbosityError Verbosity = 1
	// VerbosityWarn also logs warnings
	VerbosityWarn Verbosity = 2
	// VerbosityInfo also logs informational messages
	VerbosityInfo Verbosity = 3
	// VerbosityDebug is most verbose
	VerbosityDebug Verbosity = 4
)

// ToApexLogLevel gets the Apex log level from Verbosity
func ToApexLogLevel(v Verbosity) log.Level {
	if v >= VerbosityDebug {
		return log.DebugLevel
	}

	switch v {
	case VerbosityInfo:
		return log.InfoLevel
	case VerbosityWarn:
		return log.WarnLevel
	case VerbosityError:
		return log.ErrorLevel
	}

	return log.FatalLevel
}
