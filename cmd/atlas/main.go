package main

import (
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/yosebyte/x/log"
)

var (
	logger  *log.Logger
	version = "dev"
)

func init() {
	logger = log.NewLogger(log.Info, true)
}

func main() {
	stop := setupSignalHandler()
	parsedURL := parseArgs(os.Args)
	setLogLevel(parsedURL.Query().Get("log"))
	executeCore(parsedURL, stop)
}

func setupSignalHandler() chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	return stop
}

func parseArgs(args []string) *url.URL {
	if len(args) < 2 {
		showExitInfo()
	}
	parsedURL, err := url.Parse(args[1])
	if err != nil {
		logger.Fatal("URL parse error: %v", err)
		os.Exit(1)
	}
	return parsedURL
}

func setLogLevel(level string) {
	switch level {
	case "debug":
		logger.SetLogLevel(log.Debug)
		logger.Debug("Log level: DEBUG")
	case "info":
		logger.SetLogLevel(log.Info)
		logger.Info("Log level: INFO")
	case "warn":
		logger.SetLogLevel(log.Warn)
		logger.Warn("Log level: WARN")
	case "error":
		logger.SetLogLevel(log.Error)
		logger.Error("Log level: ERROR")
	case "fatal":
		logger.SetLogLevel(log.Fatal)
		logger.Fatal("Log level: FATAL")
	default:
		logger.SetLogLevel(log.Info)
		logger.Info("Default level: INFO")
	}
}

func showExitInfo() {
	logger.SetLogLevel(log.Info)
	logger.Info(`Version: %v %v/%v

Usage: atlas <core_mode>://<server_addr>#<access_addr>?<log=level>
`, version, runtime.GOOS, runtime.GOARCH)
	os.Exit(1)
}
