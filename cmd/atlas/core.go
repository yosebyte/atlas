package main

import (
	"crypto/tls"
	"net/url"
	"os"
	"time"

	"github.com/yosebyte/passport/pkg/log"
)

func coreSelect(parsedURL *url.URL, tlsConfig *tls.Config) {
	switch parsedURL.Scheme {
	case "server":
		runServer(parsedURL, tlsConfig)
	case "client":
		runClient(parsedURL)
	default:
		helpInfo()
		os.Exit(1)
	}
}

func runServer(parsedURL *url.URL, tlsConfig *tls.Config) {
	log.Info("Server started: %v", parsedURL.String())
	for {
		if err := server(parsedURL, tlsConfig); err != nil {
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
