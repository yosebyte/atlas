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
		runServer(parsedURL)
	case "client":
		runClient(parsedURL)
	default:
		log.Error("Invalid scheme: %v", parsedURL.Scheme)
		helpInfo()
		os.Exit(1)
	}
}

func runServer(parsedURL *url.URL) {
	log.Info("Server started: %v", parsedURL.String())
	server := internal.NewServer(parsedURL)
	if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Error("Server error: %v", err)
		}
}

func runClient(parsedURL *url.URL) {
	log.Info("Client started: %v", parsedURL.String())
	client := internal.NewClient(parsedURL)
	if err := client.ListenAndServe(); err != nil {
			log.Error("Client error: %v", err)
		}
}
