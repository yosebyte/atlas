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
	version = "dev"
	logger  = log.NewLogger(log.Info, true)
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	rawURL := os.Args[1]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		logger.Fatal("Unable to parse raw URL: %v", err)
	}
	if parsedURL.Query().Has("debug") {
		logger.SetLogLevel(log.Debug)
		logger.Debug("Debug logging enabled")
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	coreSelect(parsedURL, stop)
}

func usage() {
	logger.Fatal(`Version: %v %v/%v

Usage: atlas <core_mode>://<server_addr>#<access_addr>(?debug)
`, version, runtime.GOOS, runtime.GOARCH)
}
