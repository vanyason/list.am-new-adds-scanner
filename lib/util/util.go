package util

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Setup logger to write logs both to the file and stdout
// Logger will create file at path : log/[DD-MM-YY]/[appname].log
// For example : log/09-03-23/scratcher.log
func SetupLogging(appName string) {
	const permissions = 0774
	const osFlag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	const timeFormat = "01-02-06"

	curDate := time.Now().Format(timeFormat)
	folderPath := filepath.Join("log", curDate)
	filepath := filepath.Join(folderPath, appName+".log")

	if err := os.MkdirAll(folderPath, permissions); err != nil {
		log.Fatalf("can not setup logger (creating dir) : %s", err)
	}

	logFile, err := os.OpenFile(filepath, osFlag, permissions)
	if err != nil {
		log.Fatalf("can not setup logger (creating file) : %s", err)
	}

	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
}

// Measures time (in seconds) inside the execution scope and logs it using log lib
// Prints message in format : "message 0.1212 sec"
// Usage :  call `defer measureTimeSeconds()()“ in the beginning of a code section you want to measure
func LogSecondsPass(message string) func() {
	start := time.Now()
	return func() { log.Printf("%s %f sec", message, time.Since(start).Seconds()) }
}
