package log

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	Log.SetReportCaller(true)
	// Create a new logger instance
	logger := logrus.New()

	// Set log level (optional, defaults to InfoLevel)
	logger.SetLevel(logrus.DebugLevel)

	// Set log output (optional, defaults to os.Stdout)
	loggerFile, err := os.OpenFile("my_log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o660)
	if err != nil {
		log.Fatal(err)
	}
	logger.SetOutput(loggerFile)

	// Custom formatter (optional)
	formatter := &logrus.JSONFormatter{}
	logger.SetFormatter(formatter)

	// Log messages with different levels
	logger.Trace("Trace-level message.")
	logger.Debug("Debug-level message.")
	logger.Info("Info-level message.")
	logger.Warn("Warning-level message.")
	logger.Error("Error-level message.")

	// Close the log file gracefully
	defer loggerFile.Close() // Log to standard output by default
}
