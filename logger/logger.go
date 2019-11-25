// Package logger is an extension of the golang log package that adds support for log levels.
// Import the package, then run logger.SetLogger, providing the if the logs should be output to a log file,
// and the log level. The public constants can be used to set the log level. 0 is no logs, and 5 or higher
// is trace level logging.
// To log at a particular level, use logger.DEBUG.Printf("message") for example.
// Log levels available are ERROR, WARN, INFO, DEBUG, and TRACE.
// The package does not provide any log rotation and will overwrite the specified log file.
// log files will be placed at "./stravacommute.log". This is hardcoded.
package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// TODO: make the logfile configurable
const logFile = "./stravacommute.log"

var (
	// TRACE used for logging trace messages
	TRACE *log.Logger
	// DEBUG used for logging debug messages
	DEBUG *log.Logger
	// INFO used for logging info messages
	INFO *log.Logger
	// WARN used for logging warning messages
	WARN *log.Logger
	// ERROR used for logging error messages
	ERROR *log.Logger
)

const (
	// NoLogLevel is the integer value that means no logs are output when passed into SetLogging
	NoLogLevel = 0
	// ErrorLevel is the integer value that means only error logs are output when passed into SetLogging
	ErrorLevel = 1
	// WarnLevel is the integer value that means warn and more severe logs are output when passed into SetLogging
	WarnLevel = 2
	// InfoLevel is the integer value that means info and more severe logs are output when passed into SetLogging
	InfoLevel = 3
	// DebugLevel is the integer value that means debug and more severe logs are output when passed into SetLogging
	DebugLevel = 4
	// TraceLevel is the integer value that means trace and more severe logs are output when passed into SetLogging
	TraceLevel = 5
)

// SetLogging sets up everything needed for logging. It takes in logToFile, which if true will log
// to ./stravacommute.log, otherwise it will log to StdOut. level sets the log level. It's value
// can be one of the values of NoLogLevel, ErrorLevel, WarnLevel, InfoLevel, DebugLevel, or TraceLevel.
// The log levels are incremental, that is TRACE_LEVEL would include all higher logs, ERROR would include
// only ERROR.
// TODO: could have log rotation
func SetLogging(logToFile bool, level int) {
	dump := ioutil.Discard
	var f io.Writer
	var err error
	if logToFile {
		f, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic("Unable to open or create " + logFile + " for logging.")
		}
	} else {
		f = os.Stdout
	}
	traceOut, debugOut, infoOut, warnOut, errorOut := dump, dump, dump, dump, dump

	if level >= ErrorLevel {
		errorOut = f
	}
	if level >= WarnLevel {
		warnOut = f
	}
	if level >= InfoLevel {
		infoOut = f
	}
	if level >= DebugLevel {
		debugOut = f
	}
	if level >= TraceLevel {
		traceOut = f
	}

	TRACE = log.New(traceOut, "TRACE ", log.Ldate|log.Ltime)
	DEBUG = log.New(debugOut, "DEBUG ", log.Ldate|log.Ltime)
	INFO = log.New(infoOut, "INFO  ", log.Ldate|log.Ltime)
	WARN = log.New(warnOut, "WARN  ", log.Ldate|log.Ltime)
	ERROR = log.New(errorOut, "ERROR ", log.Ldate|log.Ltime)
}
