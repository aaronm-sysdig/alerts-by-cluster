package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strings"
)

type customFormatter struct {
	logrus.TextFormatter
}

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

func init() {

	// Define flags using pflag
	pflag.String("config", "", "Path to the config file")
	pflag.String("name", "default", "Name of the application")
	pflag.Int("port", 8080, "Port to run the application on")

	// Parse the flags
	pflag.Parse()

	// Bind pflag to viper
	viper.BindPFlags(pflag.CommandLine)
}

func initConfig() {
	// Read the config file
	configFile := viper.GetString("config")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	// Read environment variables
	viper.AutomaticEnv()

	// Read the configuration
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
	}
}

func main() {
	// Set the formatter to text
	logger := logrus.New()
	logger.SetFormatter(&customFormatter{logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	}})
	logger.SetReportCaller(true) // Enables reporting of file, function, and line number
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	// Initialize the configuration
	initConfig()

	// Access configuration values
	name := viper.GetString("name")
	port := viper.GetInt("port")
	exclude := viper.GetStringSlice("exclude")

	// Use the configuration values
	logger.Printf("Running %s on port %d\n", name, port)
	logger.Debugf("main:: %+v", exclude)
}
