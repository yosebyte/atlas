package main

import (
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/yosebyte/passport/pkg/conn"
	"github.com/yosebyte/passport/pkg/log"
	"github.com/yosebyte/passport/pkg/tls"
)

func server(parsedURL *url.URL) error {
	listenAddr := parsedURL.Host
	tlsConfig, err := acme(listenAddr)
	if err != nil {
		log.Error("Unable to obtain TLS config: %v", err)
		if tlsConfig, err = tls.NewTLSconfig(listenAddr); err != nil {
			log.Fatal("Unable to generate TLS config: %v", err)
		}
	}
	server := &http.Server{
		Addr:      listenAddr,
		TLSConfig: tlsConfig,
		ErrorLog:  log.NewLogger(),
		Handler:   http.HandlerFunc(handleServerRequest),
	}
	return server.ListenAndServeTLS("", "")
}

func handleServerRequest(w http.ResponseWriter, r *http.Request) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	targetConn, err := net.Dial("tcp", r.URL.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		clientConn.Close()
		return
	}
	if err := r.Write(targetConn); err != nil {
		log.Error("Unable to write request: %v", err)
		clientConn.Close()
		targetConn.Close()
		return
	}
	if err := conn.DataExchange(clientConn, targetConn); err != nil {
		if err == io.EOF {
			log.Info("Connection closed successfully: %v", err)
		} else {
			log.Warn("Connection closed unexpectedly: %v", err)
		}
	}
}
