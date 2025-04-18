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
	coreDispatch(parsedURL, getSignalChan())
}

func getSignalChan() chan os.Signal {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	return signalChan
}

func getParsedURL(args []string) *url.URL {
	if len(args) < 2 {
		getExitInfo()
	}
	parsedURL, err := url.Parse(args[1])
	if err != nil {
		logger.Fatal("Parse failed: %v", err)
		getExitInfo()
	}
	return parsedURL
}

func initLogLevel(level string) {
	switch level {
	case "debug":
		logger.SetLogLevel(log.Debug)
		logger.Debug("Init log level: DEBUG")
	case "warn":
		logger.SetLogLevel(log.Warn)
		logger.Warn("Init log level: WARN")
	case "error":
		logger.SetLogLevel(log.Error)
		logger.Error("Init log level: ERROR")
	case "fatal":
		logger.SetLogLevel(log.Fatal)
		logger.Fatal("Init log level: FATAL")
	default:
		logger.SetLogLevel(log.Info)
		logger.Info("Init log level: INFO")
	}
}

func getExitInfo() {
	logger.SetLogLevel(log.Info)
	logger.Info(`Version: %v %v/%v

Usage: 
    atlas <core_mode>://<server_addr>/<access_addr>?<log=level>

Examples:
    # Run as server
    atlas server://example.org/127.0.0.1:1?log=debug

    # Run as client
    atlas client://example.org:10101/127.0.0.1:10101
`, version, runtime.GOOS, runtime.GOARCH)
	os.Exit(1)
}
