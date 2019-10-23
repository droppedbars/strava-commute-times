package main

import (
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

func init() {
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic("Unable to open or create " + logFile + " for logging.")
	}
	// TODO: could make the log levels and output locations configurable
	// TODO: could have log rotation
	TRACE = log.New(ioutil.Discard, "TRACE ", log.Ldate|log.Ltime)
	DEBUG = log.New(ioutil.Discard, "DEBUG ", log.Ldate|log.Ltime)
	INFO = log.New(f, "INFO  ", log.Ldate|log.Ltime)
	WARN = log.New(f, "WARN  ", log.Ldate|log.Ltime)
	ERROR = log.New(f, "ERROR ", log.Ldate|log.Ltime)

	// TODO: need to figure out if there's a way to clean up and close the files at program end,
	//  or just let them close by the program ending
}
