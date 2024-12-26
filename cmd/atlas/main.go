package main

import (
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/yosebyte/x/log"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		helpInfo()
		os.Exit(1)
	}
	rawURL := os.Args[1]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("Unable to parse raw URL: %v", err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	coreSelect(parsedURL, stop)
}
