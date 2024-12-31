package main

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"github.com/yosebyte/atlas/internal"
)

func executeCore(parsedURL *url.URL, stop chan os.Signal) {
	switch parsedURL.Scheme {
	case "server":
		runServer(parsedURL, stop)
	case "client":
		runClient(parsedURL, stop)
	default:
		showExitInfo()
	}
}

func runServer(parsedURL *url.URL, stop chan os.Signal) {
	logger.Info("Server started: %v", parsedURL.Host)
	server := internal.NewServer(parsedURL, logger)
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error: %v", err)
		}
	}()
	<-stop
	logger.Info("Server stopping")
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Error("Server shutdown error: %v", err)
	}
	logger.Info("Server stopped")
}

func runClient(parsedURL *url.URL, stop chan os.Signal) {
	logger.Info("Client started: %v", parsedURL.Host)
	client := internal.NewClient(parsedURL, logger)
	go func() {
		if err := client.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Client error: %v", err)
		}
	}()
	<-stop
	logger.Info("Client stopping")
	if err := client.Shutdown(context.Background()); err != nil {
		logger.Error("Client shutdown error: %v", err)
	}
	logger.Info("Client stopped")
}
