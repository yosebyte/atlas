package main

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/yosebyte/atlas/internal"
	"github.com/yosebyte/x/tls"
	"golang.org/x/crypto/acme/autocert"
)

func coreDispatch(parsedURL *url.URL, signalChan chan os.Signal) {
	switch parsedURL.Scheme {
	case "server":
		runServer(parsedURL, signalChan)
	case "client":
		runClient(parsedURL, signalChan)
	default:
		logger.Fatal("Invalid scheme: %v", parsedURL.Scheme)
		getExitInfo()
	}
}

func runServer(parsedURL *url.URL, signalChan chan os.Signal) {
	if parsedURL.Hostname() != "" && net.ParseIP(parsedURL.Hostname()) == nil {
		logger.Info("Apply autocert: %v", parsedURL.Hostname())
		manager := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("autocert"),
			HostPolicy: autocert.HostWhitelist(parsedURL.Hostname()),
		}
		autocertSvr := &http.Server{
			Addr:    ":80",
			Handler: manager.HTTPHandler(nil),
		}
		internalSvr := internal.NewServer(parsedURL, manager.TLSConfig(), logger)
		go func() {
			logger.Debug("Autocert started: %v", autocertSvr.Addr)
			if err := autocertSvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("Autocert error: %v", err)
			}
		}()
		go func() {
			logger.Info("Server started: %v", parsedURL)
			logger.Info("Access address: %v", internalSvr.Addr)
			if err := internalSvr.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				logger.Error("Server error: %v", err)
			}
		}()
		<-signalChan
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go func() {
			logger.Debug("Autocert shutting down")
			if err := autocertSvr.Shutdown(ctx); err != nil {
				logger.Error("Autocert shutdown error: %v", err)
			}
			logger.Debug("Autocert shutdown complete")
		}()
		go func() {
			logger.Info("Server shutting down")
			if err := internalSvr.Shutdown(ctx); err != nil {
				logger.Error("Server shutdown error: %v", err)
			}
			logger.Info("Server shutdown complete")
		}()
	} else {
		logger.Info("Apply RAM cert: %v", version)
		tlsConfig, err := tls.GenerateTLSConfig("yosebyte/atlas:" + version)
		if err != nil {
			logger.Fatal("Generate failed: %v", err)
			return
		}
		server := internal.NewServer(parsedURL, tlsConfig, logger)
		go func() {
			logger.Info("Server started: %v", parsedURL)
			logger.Info("Access address: %v", server.Addr)
			if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				logger.Error("Server error: %v", err)
			}
		}()
		<-signalChan
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		logger.Info("Server shutting down")
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Server shutdown error: %v", err)
		}
		logger.Info("Server shutdown complete")
	}
}

func runClient(parsedURL *url.URL, signalChan chan os.Signal) {
	client := internal.NewClient(parsedURL, logger)
	go func() {
		logger.Info("Client started: %v", parsedURL)
		logger.Info("Access address: %v", client.Addr)
		if err := client.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Client error: %v", err)
		}
	}()
	<-signalChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger.Info("Client shutting down")
	if err := client.Shutdown(ctx); err != nil {
		logger.Error("Client shutdown error: %v", err)
	}
	logger.Info("Client shutdown complete")
}
