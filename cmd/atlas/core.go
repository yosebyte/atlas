package main

import (
	"net/url"
	"os"

	"github.com/yosebyte/atlas/internal"
	"github.com/yosebyte/x/log"
)

func coreSelect(parsedURL *url.URL) {
	switch parsedURL.Scheme {
	case "server":
		server(parsedURL)
	case "client":
		client(parsedURL)
	default:
		log.Error("Invalid scheme: %v", parsedURL.Scheme)
		helpInfo()
		os.Exit(1)
	}
}

func server(parsedURL *url.URL) {
	log.Info("Server started: %v", parsedURL.String())
	if err := internal.runServer(parsedURL); err != nil {
			log.Error("Server error: %v", err)
		}
}

func client(parsedURL *url.URL) {
	log.Info("Client started: %v", parsedURL.String())
	if err := internal.runClient(parsedURL); err != nil {
			log.Error("Client error: %v", err)
		}
}
