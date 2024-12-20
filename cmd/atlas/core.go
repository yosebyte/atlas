package main

import (
	"net/url"
	"os"
	"time"

	"github.com/yosebyte/passport/pkg/log"
)

func coreSelect(parsedURL *url.URL) {
	switch parsedURL.Scheme {
	case "server":
		runServer(parsedURL)
	case "client":
		runClient(parsedURL)
	default:
		helpInfo()
		os.Exit(1)
	}
}

func runServer(parsedURL *url.URL) {
	log.Info("Server started: %v", parsedURL.String())
	for {
		if err := server(parsedURL); err != nil {
			log.Error("Server error: %v", err)
			time.Sleep(1 * time.Second)
			log.Info("Server restarted")
		}
	}
}

func runClient(parsedURL *url.URL) {
	log.Info("Client started: %v", parsedURL.String())
	for {
		if err := client(parsedURL); err != nil {
			log.Error("Client error: %v", err)
			time.Sleep(1 * time.Second)
			log.Info("Client restarted")
		}
	}
}
