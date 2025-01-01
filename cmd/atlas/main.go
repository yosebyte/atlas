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
	logger  = log.NewLogger(log.Info, true)
	version = "dev"
)

func main() {
	parsedURL := getParsedURL(os.Args)
	initLogLevel(parsedURL.Query().Get("log"))
	coreManagement(parsedURL, getStopSignal())
}

func getStopSignal() chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	return stop
}

func getParsedURL(args []string) *url.URL {
	if len(args) < 2 {
		getExitInfo()
	}
	parsedURL, err := url.Parse(args[1])
	if err != nil {
		logger.Fatal("URL parse: %v", err)
		os.Exit(1)
	}
	return parsedURL
}

func initLogLevel(level string) {
	switch level {
	case "debug":
		logger.SetLogLevel(log.Debug)
		logger.Debug("Log level init: DEBUG")
	case "warn":
		logger.SetLogLevel(log.Warn)
		logger.Warn("Log level init: WARN")
	case "error":
		logger.SetLogLevel(log.Error)
		logger.Error("Log level init: ERROR")
	case "fatal":
		logger.SetLogLevel(log.Fatal)
		logger.Fatal("Log level init: FATAL")
	default:
		logger.SetLogLevel(log.Info)
		logger.Info("Log level init: INFO")
	}
}

func getExitInfo() {
	logger.SetLogLevel(log.Info)
	logger.Info(`Version: %v %v/%v

Usage: atlas <core_mode>://<server_addr>#<access_addr>?<log=level>
`, version, runtime.GOOS, runtime.GOARCH)
	os.Exit(1)
}
