package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

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
	noLogLevel = 0
	errorLevel = 2
	warnLevel  = 4
	infoLevel  = 8
	debugLevel = 16
	traceLevel = 32
)

// setLogging sets up everything needed for logging. It takes in logToFile, which if true will log
// to ./stravacommute.log, otherwise it will log to StdOut. level sets the log level. It's value
// can be one of the LogLevel struct values. The log levels are incremental, that is TRACE_LEVEL
// would include all higher logs, ERROR would include only ERROR.
// TODO: could have log rotation
func setLogging(logToFile bool, level int) {
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

	if level >= errorLevel {
		errorOut = f
	}
	if level >= warnLevel {
		warnOut = f
	}
	if level >= infoLevel {
		infoOut = f
	}
	if level >= debugLevel {
		debugOut = f
	}
	if level >= traceLevel {
		traceOut = f
	}

	TRACE = log.New(traceOut, "TRACE ", log.Ldate|log.Ltime)
	DEBUG = log.New(debugOut, "DEBUG ", log.Ldate|log.Ltime)
	INFO = log.New(infoOut, "INFO  ", log.Ldate|log.Ltime)
	WARN = log.New(warnOut, "WARN  ", log.Ldate|log.Ltime)
	ERROR = log.New(errorOut, "ERROR ", log.Ldate|log.Ltime)
}
