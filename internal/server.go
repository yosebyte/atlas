package internal

import (
	"net"
	"net/http"
	"net/url"

	"github.com/yosebyte/x/io"
	"github.com/yosebyte/x/log"
	"github.com/yosebyte/x/tls"
)

func NewServer(parsedURL *url.URL) *http.Server {
	tlsConfig, err := tls.NewTLSconfig(getagentID())
	if err != nil {
		log.Fatal("Unable to generate TLS config: %v", err)
	}
	return &http.Server{
		Addr:      parsedURL.Host,
		ErrorLog:  log.NewLogger(),
		Handler:   http.HandlerFunc(handleServerRequest),
		TLSConfig: tlsConfig,
	}
}

func handleServerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		log.Debug("User-Agent: %v", r.Header.Get("User-Agent"))
		if r.Header.Get("User-Agent") != getagentID() {
			http.Error(w, "Connection Established", http.StatusOK)
		}
		clientConn, err := hijackConnection(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Unable to hijack connection: %v", err)
			return
		}
		log.Debug("Client connected: %v", clientConn.RemoteAddr())
		defer func() {
			if clientConn != nil {
				clientConn.Close()
			}
		}()
		targetConn, err := net.Dial("tcp", r.URL.Host)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			log.Error("Unable to dial target: %v", err)
			return
		}
		log.Debug("Target connected: %v", targetConn.RemoteAddr())
		defer func() {
			if targetConn != nil {
				targetConn.Close()
			}
		}()
		log.Debug("Connection established: %v <-> %v", clientConn.RemoteAddr(), targetConn.RemoteAddr())
		if err := io.DataExchange(clientConn, targetConn); err != nil {
			log.Debug("Connection closed: %v", err)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed: %v/%v", r.RemoteAddr, r.Method)
		return
	}
}
