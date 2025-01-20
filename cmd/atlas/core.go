package main

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/yosebyte/atlas/internal"
	"github.com/yosebyte/x/tls"
	"golang.org/x/crypto/acme/autocert"
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
	var server *http.Server
	if parsedURL.Hostname() != "" && net.ParseIP(parsedURL.Hostname()) == nil {
		manager := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("autocert"),
			HostPolicy: autocert.HostWhitelist(parsedURL.Hostname()),
		}
		tlsConfig := manager.TLSConfig()
		server = internal.NewServer(parsedURL, tlsConfig, logger)
		logger.Debug("Using autocert for %v", parsedURL.Hostname())
	} else {
		tlsConfig, err := tls.NewTLSconfig("yosebyte/atlas:" + version)
		if err != nil {
			logger.Fatal("Unable to generate TLS config: %v", err)
			getExitInfo()
		}
		server = internal.NewServer(parsedURL, tlsConfig, logger)
		logger.Debug("Using self-signed certificate")
	}
	go func() {
		logger.Info("Server started: %v", parsedURL.String())
		logger.Info("Access address: %v", server.Addr)
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
