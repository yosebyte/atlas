package main

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"github.com/yosebyte/atlas/internal"
	"github.com/yosebyte/x/tls"
)

func coreDispatch(parsedURL *url.URL, stop chan os.Signal) {
	switch parsedURL.Scheme {
	case "server":
		runServer(parsedURL, stop)
	case "client":
		runClient(parsedURL, stop)
	default:
		logger.Fatal("Invalid scheme: %v", parsedURL.Scheme)
		getExitInfo()
	}
}

func runServer(parsedURL *url.URL, stop chan os.Signal) {
	tlsConfig, err := tls.NewTLSconfig("yosebyte/atlas:" + version)
	if err != nil {
		logger.Error("Unable to generate TLS config: %v", err)
	}
	server := internal.NewServer(parsedURL, tlsConfig, logger)
	go func() {
		logger.Info("Server started: %v", parsedURL.String())
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
	client := internal.NewClient(parsedURL, logger)
	go func() {
		logger.Info("Client started: %v", parsedURL.String())
		logger.Info("Access address: %v", client.Addr)
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
