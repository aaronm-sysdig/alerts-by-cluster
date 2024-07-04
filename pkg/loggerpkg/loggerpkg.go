package loggerpkg

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

// customFormatter extends logrus.TextFormatter
type customFormatter struct {
	logrus.TextFormatter
}

// initLogging initializes the logger with custom settings
func initLogging() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&customFormatter{logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	}})
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	return logger
}

// Format formats the log entry according to the customFormatter settings
func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Ensure log level is displayed with a fixed width of 5 characters
	levelText := fmt.Sprintf("%-6s", strings.ToUpper(entry.Level.String()))

	// Simplify and format the function name and line number directly after the function name
	functionName := runtime.FuncForPC(entry.Caller.PC).Name()
	lastSlash := strings.LastIndex(functionName, "/")
	shortFunctionName := functionName[lastSlash+1:] // Default to full string post-last slash if no parenthesis found

	// Try to find the parenthesis and refine shortFunctionName
	if lastSlash != -1 {
		firstParen := strings.Index(functionName[lastSlash:], "(") + lastSlash
		if firstParen > lastSlash {
			shortFunctionName = functionName[lastSlash+1 : firstParen]
		}
	}

	formattedCaller := fmt.Sprintf("%s:%d", shortFunctionName, entry.Caller.Line) // Combine function name and line

	// Right-align the caller info to ensure that it occupies a fixed width
	rightAlignedCaller := fmt.Sprintf("%-40s", formattedCaller)

	// Create the formatted log entry
	logMessage := fmt.Sprintf("%s[%s] %s %s\n", levelText, entry.Time.Format("2006-01-02 15:04:05"), rightAlignedCaller, entry.Message)

	return []byte(logMessage), nil
}

// GetLogger returns the singleton loggerpkg instance
func GetLogger() *logrus.Logger {
	once.Do(func() {
		logger = initLogging()
	})
	return logger
}
