package main

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"github.com/yosebyte/atlas/internal"
	"github.com/yosebyte/x/log"
)

func coreSelect(parsedURL *url.URL, stop chan os.Signal) {
	switch parsedURL.Scheme {
	case "server":
		runServer(parsedURL, stop)
	case "client":
		runClient(parsedURL, stop)
	default:
		log.Error("Invalid scheme: %v", parsedURL.Scheme)
		helpInfo()
		os.Exit(1)
	}
}

func runServer(parsedURL *url.URL, stop chan os.Signal) {
	log.Info("Server started: %v", parsedURL.String())
	server := internal.NewServer(parsedURL)
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Error("Server error: %v", err)
		}
	}()
	<-stop
	log.Info("Server stopping")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error("Server shutdown error: %v", err)
	}
	log.Info("Server stopped")
}

func runClient(parsedURL *url.URL, stop chan os.Signal) {
	log.Info("Client started: %v", parsedURL.String())
	client := internal.NewClient(parsedURL)
	go func() {
		if err := client.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Client error: %v", err)
		}
	}()
	<-stop
	log.Info("Client stopping")
	if err := client.Shutdown(context.Background()); err != nil {
		log.Error("Client shutdown error: %v", err)
	}
	log.Info("Client stopped")
}
